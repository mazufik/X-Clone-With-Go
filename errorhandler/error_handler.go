package errorhandler

import (
	"GO-SOCMED/dto"
	"GO-SOCMED/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var statuscode int

	switch err.(type) {
	case *NotFoundError:
		statuscode = http.StatusNotFound
	case *BadRequestError:
		statuscode = http.StatusBadRequest
	case *InternalServerError:
		statuscode = http.StatusInternalServerError
	case *UnauthorizedError:
		statuscode = http.StatusUnauthorized
	}

	response := helper.Response(dto.ResponseParams{
		StatusCode: statuscode,
		Message:    err.Error(),
	})

	c.JSON(statuscode, response)
}
