package gateway

import (
	"log"
	"net/http"
	"sn/gateway/internal/crypto"
	"time"

	"github.com/gin-gonic/gin"

	gen "sn/gateway/generated"
	"sn/gateway/internal/usecase"
)

func (*Server) RegisterUser(ctx *gin.Context) {
	body := gen.UserRegistration{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gen.Error{
			Message: "Invalid request body",
			Code:    "bad_request",
		})
		return
	}

	profile, err := usecase.RegisterUser(&body)
	if err != nil {
		log.Printf("failed to register user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to register user",
			Code:    "internal_server_error",
		})
		return
	}
	if profile.GetId() == "" {
		log.Printf("failed to get id of newly registered user")
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to register user",
			Code:    "internal_server_error",
		})
	}

	token, err := crypto.CreateToken(profile.GetId(), time.Now().Add(3600*time.Second))
	if err != nil {
		log.Printf("failed to create token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gen.Error{
			Message: "Failed to create token",
			Code:    "internal_server_error",
		})
		return
	}

	ctx.SetCookie("auth_token", token, 3600, "/", "localhost", false, true)
	ctx.Status(http.StatusCreated)
}
