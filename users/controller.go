package users

import (
  "net/http"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "time"
  
  "github.com/gin-gonic/gin"
)

type controller struct {
  rg *gin.RouterGroup
  useCase UseCase
}

func (c *controller) insertHandler(ctx *gin.Context) {
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

  user, err := c.useCase.Register(payload, path)
  if err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  createdAt := user.CreatedAt.Format(time.RFC850)
  common.SendCreatedResponse(ctx, user, createdAt, "Register Successfully")
}

func (c *controller) Route() {
  c.rg.POST("/auth/register/merchants", c.insertHandler)
  c.rg.POST("/auth/register", c.insertHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase) *controller {
  return &controller{rg, useCase}
}
