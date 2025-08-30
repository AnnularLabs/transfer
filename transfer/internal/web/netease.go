package web

import (
	"net/http"
	"strconv"
	"transfer/internal/service"

	"github.com/gin-gonic/gin"
)

var _ handler = (*NetEaseHandler)(nil)

// NetEaseHandler 网易云音乐
type NetEaseHandler struct {
	svc service.NeteaseService
}

func NewNetEaseHandler(svc service.NeteaseService) *NetEaseHandler {
	return &NetEaseHandler{
		svc: svc,
	}
}

func (n *NetEaseHandler) RegisterRoutes(server *gin.Engine) {
	ng := server.Group("/netease")
	ng.GET("/playlist", n.GetPlaylist)
}

func (n *NetEaseHandler) GetPlaylist(ctx *gin.Context) {
	nIdStr := ctx.Query("id")

	nId, _ := strconv.ParseInt(nIdStr, 10, 64)
	playlist, err := n.svc.GetPlaylist(ctx.Request.Context(), nId)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}
