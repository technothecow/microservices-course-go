package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gen "sn/gateway/generated"
)

func (*Server) Ping(ctx *gin.Context) {
	resp := gen.PingResponse{
		Message: "pong",
	}

	ctx.JSON(http.StatusOK, resp)
}
