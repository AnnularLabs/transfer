package api

import (
	"net/http"
	"transfer/internal/netease"

	"github.com/gin-gonic/gin"
)

type NeteaseRequest struct {
	PlaylistId string `json:"PlaylistId"`
}

type resultResponse struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func (server *Server) Netease(gtx *gin.Context) {
	var req NeteaseRequest
	result, _ := gtx.GetQuery("PlaylistId")
	if result == "" {
		gtx.JSON(http.StatusBadRequest, "PlaylistId is required")
		return
	}
	req.PlaylistId = result
	musicList, err := netease.NeteaseService.GetPlayListMusic(gtx.Request.Context(), req.PlaylistId)
	if err != nil {
		gtx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := &resultResponse{
		Msg:  "success",
		Data: musicList,
	}
	gtx.JSON(http.StatusOK, res)
}
