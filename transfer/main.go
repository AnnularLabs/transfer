package main

import (
	"transfer/internal/service"
	"transfer/internal/service/oauth2"
	"transfer/internal/service/session"
	"transfer/internal/web"
	"transfer/internal/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/zmb3/spotify"
)

func main() {
	server := initWeb()
	server.Run(":8081")
}

func initSpotifyClient() spotify.Client {
	// oauth2 - 保持原有的应用级别客户端（用于公开API调用）
	return oauth2.NewSpotifyAuth("your-client-id", "your-client-secret")
}
func initWeb() *gin.Engine {
	// 1. 初始化组件
	tokenManager := oauth2.NewMemoryTokenManager()
	sessionManager := session.NewMemorySessionManager()
	oauthService := oauth2.NewSpotifyOAuth(
		"your-client-id",
		"your-client-secret",
		"your-callback-url",
		tokenManager,
	)

	// 2. 初始化服务和处理器
	nsv := service.NewNeteaseService()
	neteaseHdl := web.NewNetEaseHandler(nsv)

	ssv := service.NewSpotifyService(initSpotifyClient())
	spotifyHdl := web.NewSpotifyHandler(ssv, tokenManager, sessionManager, oauthService)

	userHdl := web.NewUserHandler(oauthService, tokenManager, sessionManager)

	// 3. 配置服务器
	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// 4. 注册路由
	neteaseHdl.RegisterRoutes(server)
	spotifyHdl.RegisterRoutes(server)
	userHdl.RegisterRoutes(server)

	return server
}
