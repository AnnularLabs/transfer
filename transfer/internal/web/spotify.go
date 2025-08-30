// 这是对现有 /Users/merrick/t/transfer/transfer/internal/web/spotify.go 的修改建议
package web

import (
	"net/http"
	"transfer/internal/domain"
	"transfer/internal/service"
	"transfer/internal/service/oauth2"
	"transfer/internal/service/session"
	"transfer/internal/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

var _ handler = (*SpotifyHandler)(nil)

type SpotifyHandler struct {
	svc            service.SpotifyService
	tokenManager   oauth2.TokenManager
	sessionManager session.SessionManager
	oauthService   oauth2.SpotifyOAuthService
}

func NewSpotifyHandler(svc service.SpotifyService, tokenManager oauth2.TokenManager, sessionManager session.SessionManager, oauthService oauth2.SpotifyOAuthService) *SpotifyHandler {
	return &SpotifyHandler{
		svc:            svc,
		tokenManager:   tokenManager,
		sessionManager: sessionManager,
		oauthService:   oauthService,
	}
}

func (s *SpotifyHandler) RegisterRoutes(server *gin.Engine) {
	sg := server.Group("/spotify")

	// 需要授权的路由
	authRequired := sg.Group("")
	authRequired.Use(middleware.RequireSpotifyAuth(s.tokenManager, s.sessionManager, s.oauthService))
	{
		authRequired.GET("/me", s.Me)
		authRequired.GET("/playlists", s.GetPlaylistsForUser)
		authRequired.POST("/playlists/:id/tracks", s.AddTracksToPlaylist)
	}
}

func (s *SpotifyHandler) Me(ctx *gin.Context) {
	// 从中间件获取 Spotify 客户端
	client, exists := ctx.Get("spotify_client")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Spotify client not found"})
		return
	}

	spotifyClient := client.(spotify.Client)

	// 获取当前用户信息
	user, err := spotifyClient.CurrentUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_get_user",
			"message": "无法获取用户信息",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"display_name": user.DisplayName,
		"email":        user.Email,
		"followers":    user.Followers.Count,
		"country":      user.Country,
	})
}

func (s *SpotifyHandler) GetPlaylistsForUser(ctx *gin.Context) {
	// 从中间件获取 Spotify 客户端
	client, exists := ctx.Get("spotify_client")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Spotify client not found"})
		return
	}

	spotifyClient := client.(spotify.Client)

	// 获取用户歌单
	playlists, err := spotifyClient.CurrentUsersPlaylists()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed_to_get_playlists",
			"message": "无法获取歌单列表",
			"details": err.Error(),
		})
		return
	}

	// 转换为前端需要的格式
	result := make([]map[string]interface{}, 0, len(playlists.Playlists))
	for _, playlist := range playlists.Playlists {
		result = append(result, map[string]interface{}{
			"id":   playlist.ID.String(),
			"name": playlist.Name,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *SpotifyHandler) AddTracksToPlaylist(ctx *gin.Context) {
	playlistId := ctx.Param("id")
	var req struct {
		TrackNames []string `json:"track_names"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "请求格式错误",
			"details": err.Error(),
		})
		return
	}

	// 从中间件获取 Spotify 客户端
	client, exists := ctx.Get("spotify_client")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Spotify client not found"})
		return
	}

	spotifyClient := client.(spotify.Client)

	// 构建歌曲列表
	var tracks []domain.Track
	for _, name := range req.TrackNames {
		tracks = append(tracks, domain.Track{Title: name})
	}

	// 这里需要修改 SpotifyService 来接受已授权的客户端
	// 或者直接在这里处理迁移逻辑
	result, err := s.svc.TransferTracksWithUserClient(ctx.Request.Context(), spotifyClient, playlistId, tracks)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "transfer_failed",
			"message": "歌曲迁移失败",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "迁移完成",
		"result":  result,
	})
}
