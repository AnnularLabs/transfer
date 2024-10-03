package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type helloTransferResponse struct {
	Message string `json:"message"`
}

func (server *Server) helloTransfer(c *gin.Context) {
	var req helloTransferResponse
	req.Message = "Hello from transfer"
	c.JSON(http.StatusOK, req)
}
