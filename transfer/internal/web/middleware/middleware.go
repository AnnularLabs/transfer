package middleware

import (
	"net/http"
	"transfer/internal/service/oauth2"
	"transfer/internal/service/session"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequireSpotifyAuth Spotify 授权中间件
func RequireSpotifyAuth(tokenManager oauth2.TokenManager, sessionManager session.SessionManager, oauthService oauth2.SpotifyOAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查 session
		sessionData := sessionManager.GetSession(c)
		if sessionData == nil || !sessionData.IsAuthed {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "spotify_auth_required",
				"message": "需要 Spotify 授权",
			})
			c.Abort()
			return
		}

		// 2. 获取已授权的 Spotify 客户端
		client, err := oauthService.GetAuthenticatedClient(sessionData.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "failed_to_get_client",
				"message": "无法获取 Spotify 客户端",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 3. 将客户端和用户信息添加到上下文
		c.Set("spotify_client", client)
		c.Set("spotify_user_id", sessionData.UserID)
		c.Set("session_data", sessionData)

		c.Next()
	}
}

// SetUserID 从查询参数或表单中获取用户ID并设置到上下文的中间件
func SetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string

		// 尝试从查询参数获取
		if id := c.Query("userId"); id != "" {
			userID = id
		}

		// 尝试从表单获取
		if userID == "" {
			if id := c.PostForm("userId"); id != "" {
				userID = id
			}
		}

		// 尝试从 JSON body 获取
		if userID == "" {
			var body struct {
				UserID string `json:"userId"`
			}
			if err := c.ShouldBindJSON(&body); err == nil && body.UserID != "" {
				userID = body.UserID
			}
		}

		if userID != "" {
			c.Set("userId", userID)
		}

		c.Next()
	}
}
