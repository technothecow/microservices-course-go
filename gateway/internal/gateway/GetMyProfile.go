package gateway

import (
	"net/http"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetMyProfile(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	profile, err := usecase.GetProfile(userId)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
