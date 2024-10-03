package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() (*Server, error) {

	server := &Server{}
	server.setupRouter()
	return server, nil

}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/helloTransfer", server.helloTransfer)
	// router.POST("/song", server.song)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}