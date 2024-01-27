package products

import (
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/middlewares"
  "net/http"
  "time"
  "strconv"

  "github.com/gin-gonic/gin"
)

type controller struct {
  rg *gin.RouterGroup
  useCase UseCase
  authMiddleware middlewares.AuthMiddleware
}

func (c *controller) insertHandler(ctx *gin.Context) {
  var id = ctx.MustGet("userId").(string)
  var payload InsertProductInput
  if err := ctx.ShouldBindJSON(&payload); err != nil {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
    return
  }
  
  if err := common.ValidateInput(payload); len(err) > 0 {
    common.SendErrorResponse(ctx, http.StatusBadRequest, err)
    return
  }
  
  product, statusCode, err := c.useCase.CreateNewProduct(payload, id)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }
  
  common.SendCreatedResponse(ctx, product, time.Now().Format("Monday, 02 January 2006 15:04:05 MST"), "Create product successfully")
}

func (c *controller) getAllProductsHandler(ctx *gin.Context) {
  page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
  size, _ := strconv.Atoi(ctx.DefaultQuery("size", "5"))

  products, paging, statusCode, err := c.useCase.GetAllProducts(page, size)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  response := []any{}
  for _, val := range products {
    response = append(response, val)
  }

  common.SendPagedResponse(ctx, response, paging, "Get all products successfully")
}

func (c *controller) getMerchantProductHandler(ctx *gin.Context) {
  page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
  size, _ := strconv.Atoi(ctx.DefaultQuery("size", "5"))
  merchantId := ctx.MustGet("userId").(string)

  products, paging, err := c.useCase.GetMerchantProducts(merchantId, page, size)
  if err != nil {
    common.SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
    return
  }

  response := []any{}
  for _, val := range products {
    response = append(response, val)
  }

  common.SendPagedResponse(ctx, response, paging, "Your products")
}

func (c *controller) getProductByIDHandler(ctx *gin.Context) {
  id := ctx.Param("id")

  product, statusCode, err := c.useCase.GetProductByID(id)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  common.SendSingleResponseWithData(ctx, product, "Get product successfully")
}

func (c *controller) deleteProductHandler(ctx *gin.Context) {
  id := ctx.Param("id")
  userId := ctx.MustGet("userId").(string)

  statusCode, err := c.useCase.DeleteProductByID(id, userId)
  if err != nil {
    common.SendErrorResponse(ctx, statusCode, err.Error())
    return
  }

  common.SendDeletedResponse(ctx, "Delete product successfully")
}

func (c *controller) Route() {
  // Common endpoint
  c.rg.GET(configs.Products, c.authMiddleware.RequireToken("MERCHANT", "CUSTOMER"), c.getAllProductsHandler)
  c.rg.GET(configs.ProductsWithIDParam, c.authMiddleware.RequireToken("MERCHANT", "CUSTOMER"), c.getProductByIDHandler)
  
  // Merchant endpoint
  merchant := c.rg.Group(configs.MerchantsGroup)
  merchant.POST(configs.Products, c.authMiddleware.RequireToken("MERCHANT"), c.insertHandler)
  merchant.GET(configs.Products, c.authMiddleware.RequireToken("MERCHANT"), c.getMerchantProductHandler)
  merchant.DELETE(configs.ProductsWithIDParam, c.authMiddleware.RequireToken("MERCHANT"), c.deleteProductHandler)
}

func NewController(rg *gin.RouterGroup, useCase UseCase, authMiddleware middlewares.AuthMiddleware) *controller {
  return &controller{rg, useCase, authMiddleware}
}
