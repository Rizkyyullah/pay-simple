package common

import (
	"github.com/Rizkyyullah/pay-simple/shared/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendCreatedResponse(c *gin.Context, data any, createdAt string, message string) {
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

func SendCreatedResponseWithoutData(c *gin.Context, createdAt string, message string) {
	c.JSON(http.StatusCreated, &models.Meta{
		Status:    "Success",
		Code:      http.StatusCreated,
		Message:   message,
		CreatedAt: createdAt,
	})
}

func SendUpdatedResponse(c *gin.Context, data any, updatedAt string, message string) {
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

func SendSingleResponseWithData(c *gin.Context, data any, message string) {
	c.JSON(http.StatusOK, &models.SingleResponse{
		Meta: models.Meta{
			Status:  "Success",
			Code:    http.StatusOK,
			Message: message,
		},
		Data: data,
	})
}

func SendSingleResponseWithoutData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, &models.Meta{
		Status:  "Success",
		Code:    http.StatusOK,
		Message: message,
	})
}

func SendPagedResponse(c *gin.Context, data []any, paging models.Paging, message string) {
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

func SendUnauthorizedResponse(c *gin.Context, message any) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, &models.Meta{
		Status:  "Error",
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

func SendForbiddenResponse(c *gin.Context, message any) {
	c.AbortWithStatusJSON(http.StatusForbidden, &models.Meta{
		Status:  "Error",
		Code:    http.StatusForbidden,
		Message: message,
	})
}
