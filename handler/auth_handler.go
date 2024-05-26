package handler

import (
	"GO-SOCMED/dto"
	"GO-SOCMED/errorhandler"
	"GO-SOCMED/helper"
	"GO-SOCMED/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}

// method
func (h *authHandler) Register(c *gin.Context) {
	var register dto.RegisterRequest

	// model binding and validation
	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		// agar code dibawah tidak tereksikusi tambahkan return
		return
	}

	if err := h.service.Register(&register); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Register successfully, please login",
	})

	c.JSON(http.StatusCreated, res)
}
