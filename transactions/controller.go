package transactions

import (
  "fmt"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/middlewares"
  "net/http"
  "strconv"
  "time"

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

func (c *controller) createTransactionHandler(ctx *gin.Context) {
  var userId = ctx.MustGet("userId").(string)
  var payload CreateTransactionInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

  if err := common.ValidateInput(payload); len(err) > 0 {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err)
    return
  }

  if err := c.useCase.CreateTransaction(userId, payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  common.SendCreatedResponseWithoutData(ctx, time.Now().In(common.GetTimezone()).Format("Monday, 02 January 2006 15:04:05"), "Create transaction successfully")
}

func (c *controller) topupHandler(ctx *gin.Context) {
  var userId = ctx.MustGet("userId").(string)
  var payload TopupInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

  if err := common.ValidateInput(payload); len(err) > 0 {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err)
    return
  }

  balance, err := c.useCase.Topup(payload, userId)
  if err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  common.SendSingleResponseWithData(ctx, balance, "Topup successfully")
}

func (c *controller) transferHandler(ctx *gin.Context) {
  var senderId = ctx.MustGet("userId").(string)
  var payload TransferInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

  if err := common.ValidateInput(payload); len(err) > 0 {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err)
    return
  }

  balance, err := c.useCase.Transfer(payload, senderId)
  if err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  common.SendSingleResponseWithData(ctx, balance, "Transfer successfully")
}

func (c *controller) Route() {
  // Customer endpoint
  customer := c.rg.Group(configs.CustomersGroup)
  customer.GET(configs.CustomerTransaction, c.authMiddleware.RequireToken("CUSTOMER"), c.getTransactionsHistoryHandler)
  customer.GET(configs.CustomerTransactionWithIDParam, c.authMiddleware.RequireToken("CUSTOMER"), c.getTransactionHistoryByIDHandler)
  customer.POST(configs.CustomerTransaction, c.authMiddleware.RequireToken("CUSTOMER"), c.createTransactionHandler)
  customer.POST(fmt.Sprintf("%s%s", configs.Balance, configs.CustomerTopup), c.authMiddleware.RequireToken("CUSTOMER"), c.topupHandler)
  customer.POST(configs.Transfer, c.authMiddleware.RequireToken("CUSTOMER"), c.transferHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase, authMiddleware middlewares.AuthMiddleware) *controller {
  return &controller{rg, useCase, authMiddleware}
}
