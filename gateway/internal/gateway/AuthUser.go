package gateway

import (
	"log"
	"net/http"
	gen "sn/gateway/generated"
	"sn/gateway/internal/crypto"
	"sn/gateway/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

func (*Server) AuthUser(ctx *gin.Context) {
	body := gen.UsernameAndPassword{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Printf("failed to bind body: %v", err)
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	userId, err := usecase.AuthUser(body)
	if err != nil {
		log.Printf("failed to auth user: %v", err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	token, err := crypto.CreateToken(userId, time.Now().Add(3600*time.Second))
	if err != nil {
		log.Printf("failed to create token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to create token",
			Code:    "internal_server_error",
		})
		return
	}

	ctx.SetCookie("auth_token", token, 3600, "/", "localhost", false, true)
	ctx.Status(http.StatusOK)
}
