package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NeteaseRequest struct {
	PlaylistId string `json:"PlaylistId"`
}

type helloTransferResponse struct {
	Message string `json:"message"`
}

type resultResponse struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (server *Server) HelloTransfer(c *gin.Context) {
	var req helloTransferResponse
	req.Message = "Hello from transfer"
	c.JSON(http.StatusOK, req)
}

func (server *Server) Netease(ctx *gin.Context) {
	var req NeteaseRequest
	result, _ := ctx.GetQuery("PlaylistId")
	if result == "" {
		ctx.JSON(http.StatusBadRequest, "PlaylistId is required")
		return
	}
	req.PlaylistId = result
	musicList, err := GetPlayListMusic(req.PlaylistId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := &resultResponse{
		Msg:  "success",
		Data: musicList,
	}
	ctx.JSON(http.StatusOK, res)
}
