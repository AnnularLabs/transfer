package web

import (
	"fmt"
	"net/http"
	"time"
	"transfer/internal/service/oauth2"
	"transfer/internal/service/session"

	xoauth2 "golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	oauthService   oauth2.SpotifyOAuthService
	tokenManager   oauth2.TokenManager
	sessionManager session.SessionManager
}

func NewUserHandler(oauthService oauth2.SpotifyOAuthService, tokenManager oauth2.TokenManager, sessionManager session.SessionManager) *UserHandler {
	return &UserHandler{
		oauthService:   oauthService,
		tokenManager:   tokenManager,
		sessionManager: sessionManager,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")

	ug.GET("/auth/spotify/login", u.InitiateAuth)
	ug.GET("/auth/spotify/callback", u.HandleCallback)
	ug.POST("/auth/spotify/status", u.CheckAuthStatus)
	ug.POST("/auth/spotify/logout", u.Logout)
}

// InitiateAuth 发起 Spotify 授权
func (u *UserHandler) InitiateAuth(c *gin.Context) {
	// 1. 生成随机 state（CSRF 保护）
	state := oauth2.GenerateSecureState()

	// 2. 存储 state 到 session
	err := u.sessionManager.SetState(c, state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_set_state",
			"message": "无法设置授权状态",
		})
		return
	}

	// 3. 生成授权 URL
	authURL := u.oauthService.GetAuthURL(state)

	// 4. 重定向到 Spotify 授权页面
	c.Redirect(http.StatusFound, authURL)
}

// HandleCallback 处理 Spotify OAuth 回调
func (u *UserHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	error := c.Query("error")

	// 1. 检查是否有错误
	if error != "" {
		frontendURL := fmt.Sprintf("http://localhost:3000?auth=error&error=%s", error)
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 2. 检查必要参数
	if code == "" || state == "" {
		frontendURL := "http://localhost:3000?auth=error&error=missing_parameters"
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 3. 验证 state（CSRF 保护）
	if !u.sessionManager.ValidateState(c, state) {
		frontendURL := "http://localhost:3000?auth=error&error=invalid_state"
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 4. 交换授权码获取 token
	token, err := u.oauthService.HandleCallback(code, state)
	if err != nil {
		frontendURL := fmt.Sprintf("http://localhost:3000?auth=error&error=token_exchange_failed&details=%s", err.Error())
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 5. 使用 token 获取用户信息
	client := spotify.Authenticator{}.NewClient(&xoauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.ExpiresAt,
	})

	user, err := client.CurrentUser()
	if err != nil {
		frontendURL := fmt.Sprintf("http://localhost:3000?auth=error&error=failed_to_get_user&details=%s", err.Error())
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 6. 存储用户 token
	userID := user.ID
	token.UserID = userID
	err = u.tokenManager.StoreUserToken(userID, token)
	if err != nil {
		frontendURL := fmt.Sprintf("http://localhost:3000?auth=error&error=failed_to_store_token&details=%s", err.Error())
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 7. 设置 session
	sessionData := &session.SessionData{
		UserID:    userID,
		SpotifyID: userID,
		IsAuthed:  true,
		AuthTime:  time.Now(),
	}
	err = u.sessionManager.SetSession(c, sessionData)
	if err != nil {
		frontendURL := fmt.Sprintf("http://localhost:3000?auth=error&error=failed_to_set_session&details=%s", err.Error())
		c.Redirect(http.StatusFound, frontendURL)
		return
	}

	// 8. 重定向回前端成功页面
	frontendURL := fmt.Sprintf("http://localhost:3000?auth=success&user=%s", userID)
	c.Redirect(http.StatusFound, frontendURL)
}

// CheckAuthStatus 检查用户授权状态
func (u *UserHandler) CheckAuthStatus(c *gin.Context) {
	// 1. 获取 session
	sessionData := u.sessionManager.GetSession(c)
	if sessionData == nil || !sessionData.IsAuthed {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"message":       "用户未授权",
		})
		return
	}

	// 2. 检查 token 是否仍然有效
	isValid := u.tokenManager.IsTokenValid(sessionData.UserID)
	if !isValid {
		// 3. 尝试刷新 token
		token, err := u.tokenManager.GetUserToken(sessionData.UserID)
		if err == nil && token.RefreshToken != "" {
			newToken, err := u.oauthService.RefreshUserToken(token.RefreshToken)
			if err == nil {
				newToken.UserID = sessionData.UserID
				u.tokenManager.StoreUserToken(sessionData.UserID, newToken)
				isValid = true
			}
		}
	}

	// 4. 如果 token 刷新失败，清除 session
	if !isValid {
		u.sessionManager.DeleteSession(c)
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"message":       "Token 已过期且无法刷新",
		})
		return
	}

	// 5. 返回授权状态
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"user_id":       sessionData.SpotifyID,
		"auth_time":     sessionData.AuthTime,
		"message":       "用户已授权",
	})
}

// Logout 用户登出
func (u *UserHandler) Logout(c *gin.Context) {
	// 1. 获取 session
	sessionData := u.sessionManager.GetSession(c)
	if sessionData != nil && sessionData.IsAuthed {
		// 2. 删除用户 token
		u.tokenManager.DeleteUserToken(sessionData.UserID)

		// 3. 清除 session
		u.sessionManager.DeleteSession(c)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}
