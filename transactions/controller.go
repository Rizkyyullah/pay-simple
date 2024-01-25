package transactions

import (
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/middlewares"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
)

type controller struct {
  rg *gin.RouterGroup
  useCase UseCase
  authMiddleware middlewares.AuthMiddleware
}

func (c *controller) getTransactionsHistoryHandler(ctx *gin.Context) {
  id := ctx.MustGet("userId").(string)
  page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
  size, _ := strconv.Atoi(ctx.DefaultQuery("size", "5"))
  transactions, paging, err := c.useCase.GetTransactionsHistory(id, page, size)
  if err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }
  
  response := []any{}
  for _, val := range transactions {
    response = append(response, val)
  }

  common.SendPagedResponse(ctx, response, paging, "Your transactions history")
}

func (c *controller) getTransactionHistoryByIDHandler(ctx *gin.Context) {
  id := ctx.Param("id")
  userId := ctx.MustGet("userId").(string)
  transaction, statusCode, err := c.useCase.GetTransactionHistoryByID(id, userId)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  common.SendSingleResponseWithData(ctx, transaction, "Your transaction history")
}

func (c *controller) Route() {
  // Customer endpoint
  customer := c.rg.Group(configs.CustomersGroup)
  customer.GET(configs.CustomerTransaction, c.authMiddleware.RequireToken("CUSTOMER"), c.getTransactionsHistoryHandler)
  customer.GET(configs.CustomerTransactionWithIDParam, c.authMiddleware.RequireToken("CUSTOMER"), c.getTransactionHistoryByIDHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase, authMiddleware middlewares.AuthMiddleware) *controller {
  return &controller{rg, useCase, authMiddleware}
}
