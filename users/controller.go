package users

import (
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/middlewares"

  "github.com/gin-gonic/gin"
)

type controller struct {
  rg *gin.RouterGroup
  useCase UseCase
  authMiddleware middlewares.AuthMiddleware
}

func (c *controller) getBalanceHandler(ctx *gin.Context) {
  id := ctx.MustGet("userId").(string)
  balance, statusCode, err := c.useCase.GetBalance(id)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  common.SendSingleResponse(ctx, balance, "Your current balance")
}

func (c *controller) Route() {
  // Merchant endpoint
  merchant := c.rg.Group(configs.MerchantsGroup)
  merchant.GET(configs.Balance, c.authMiddleware.RequireToken("MERCHANT"), c.getBalanceHandler)
  
  // Customer endpoint
  customer := c.rg.Group(configs.CustomersGroup)
  customer.GET(configs.Balance, c.authMiddleware.RequireToken("CUSTOMER"), c.getBalanceHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase, authMiddleware middlewares.AuthMiddleware) *controller {
  return &controller{rg, useCase, authMiddleware}
}
