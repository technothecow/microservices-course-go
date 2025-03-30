package gateway

import (
	"github.com/gin-gonic/gin"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
)

func (*Server) EditMyProfile(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.EditProfile{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	profile, err := usecase.EditProfile(userId, &body)
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
