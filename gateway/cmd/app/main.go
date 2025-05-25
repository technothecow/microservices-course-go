package main

import (
	"log"
	"net/http"
	"sn/libraries/kafka"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/gin-middleware"
	"go.uber.org/zap"

	gen "sn/gateway/generated"
	"sn/gateway/internal/gateway"
)

func main() {
	logger := zap.NewExample()
	logger.Info("Starting server on port 50001")

	server := gateway.NewServer()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	defer kafka.CloseProducer()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	validationMiddleware, err := ginmiddleware.OapiValidatorFromYamlFile("openapi.yaml")
	if err != nil {
		panic(err)
	}
	router.Use(validationMiddleware)
	gen.RegisterHandlers(router, server)
	s := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:50001",
	}

	log.Fatal(s.ListenAndServe())
}
