package auth

import (
  "net/http"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "time"
  
  "github.com/gin-gonic/gin"
)

type controller struct {
  rg *gin.RouterGroup
  useCase UseCase
  jwtService services.JwtService
}

func (c *controller) registerHandler(ctx *gin.Context) {
  var path = ctx.Request.RequestURI
  var payload RegisterInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

  if err := common.ValidateInput(payload); len(err) > 0 {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err)
    return
  }

  user, statusCode, err := c.useCase.Register(payload, path)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  // Send userId to context to find out identity when loggingout
  ctx.Set("userId", user.ID)

  common.SendCreatedResponse(ctx, user, time.Now().In(common.GetTimezone()).Format("Monday, 02 January 2006 15:04:05 WIB"), "Register Successfully")
}

func (c *controller) loginHandler(ctx *gin.Context) {
  var path = ctx.Request.RequestURI
  var payload LoginInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }

  token, statusCode, err := c.useCase.Login(payload, path)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  // Send email to context to find out identity when loggingout
  ctx.Set("email", payload.Email)

  ctx.SetCookie("auth_cookie", token, 3600, configs.APIGroup + "/", "", false, true)
  common.SendSingleResponseWithoutData(ctx, "Login Successfully")
}

func (c *controller) logoutHandler(ctx *gin.Context) {
  biscuit, err := ctx.Cookie("auth_cookie")
  if err != nil {
    common.SendUnauthorizedResponse(ctx, "You haven't logged in to the application")
		return
  }

  claims, err := c.jwtService.ParseToken(biscuit)
  if err != nil {
    common.SendUnauthorizedResponse(ctx, "You haven't logged in to the application")
		return
  }

  claims["authorized"] = false
  ctx.SetCookie("auth_cookie", "", -1, configs.APIGroup + "/", "", false, true)

  // Send userId to context to find out identity when loggingout
  ctx.Set("userId", claims["userId"])

  common.SendSingleResponseWithoutData(ctx, "Logout Successfully")
}

func (c *controller) Route() {
  // Common endpoint
  c.rg.GET(configs.Logout, c.logoutHandler)

  // Merchant endpoint
  c.rg.POST(configs.MerchantRegister, c.registerHandler)
  c.rg.POST(configs.MerchantLogin, c.loginHandler)

  // Customer endpoint
  c.rg.POST(configs.CustomerRegister, c.registerHandler)
  c.rg.POST(configs.CustomerLogin, c.loginHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase, jwtService services.JwtService) *controller {
  return &controller{rg, useCase, jwtService}
}
