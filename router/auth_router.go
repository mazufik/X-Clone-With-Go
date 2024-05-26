package router

import (
	"GO-SOCMED/config"
	"GO-SOCMED/handler"
	"GO-SOCMED/repository"
	"GO-SOCMED/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	api.POST("/register", authHandler.Register)
}
