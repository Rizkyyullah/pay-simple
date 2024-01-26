package users

import (
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/middlewares"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/jinzhu/copier"
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

  common.SendSingleResponseWithData(ctx, balance, "Your current balance")
}

func (c *controller) getProfileHandler(ctx *gin.Context) {
  id := ctx.MustGet("userId").(string)
  user, statusCode, err := c.useCase.GetUserByID(id)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  var profile UserResponse
  if err := copier.Copy(&profile, &user); err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  common.SendSingleResponseWithData(ctx, profile, "Your profile")
}

func (c *controller) Route() {
  // Common endpoint
  c.rg.GET(configs.Profile, c.authMiddleware.RequireToken("MERCHANT", "CUSTOMER"), c.getProfileHandler)
  
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
