package common

import (
	"github.com/Rizkyyullah/pay-simple/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendCreatedResponse(c *gin.Context, data interface{}, createdAt string, message string) {
	c.JSON(http.StatusCreated, &models.SingleResponse{
		Meta: models.Meta{
			Status:    "Success",
			Code:      http.StatusCreated,
			Message:   message,
			CreatedAt: createdAt,
		},
		Data: data,
	})
}

func SendUpdatedResponse(c *gin.Context, data interface{}, updatedAt string, message string) {
	c.JSON(http.StatusCreated, &models.SingleResponse{
		Meta: models.Meta{
			Status:    "Success",
			Code:      http.StatusCreated,
			Message:   message,
			UpdatedAt: updatedAt,
		},
		Data: data,
	})
}

func SendSingleResponse(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, &models.SingleResponse{
		Meta: models.Meta{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func SendPagedResponse(c *gin.Context, data []interface{}, paging models.Paging, message string) {
	c.JSON(http.StatusOK, &models.PagedResponse{
		Meta: models.Meta{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: message,
		},
		Data:   data,
		Paging: paging,
	})
}

func SendDeletedResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, &models.Meta{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: message,
	},
	)
}

func SendErrorResponse(c *gin.Context, code int, message any) {
	c.JSON(code, &models.Meta{
		Status:  "Error",
		Code:    code,
		Message: message,
	})
}
