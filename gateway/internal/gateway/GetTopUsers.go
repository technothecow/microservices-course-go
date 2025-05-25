package gateway

import (
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (*Server) GetTopUsers(ctx *gin.Context) {
	_, err := GetUserIdFromContext(ctx)
	if err != nil {
		return
	}

	body := gen.TopUsersRequest{}
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "invalid_request_body",
		})
		return
	}

	comments, err := usecase.GetTopUsers(&body)
	if err != nil {
		log.Printf("Failed to get top Users: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to get top Users",
			Code:    "failed_to_get_top_Users",
		})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}