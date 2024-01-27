package products

import (
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "github.com/Rizkyyullah/pay-simple/users"
  "log"
  "net/http"
  
  "github.com/jinzhu/copier"
)

type S []string
type M map[string]entities.User

func (s S) contains(value string) bool {
  for _, val := range s {
    if val == value {
      return true
    }
  }

  return false
}

type UseCase interface {
  CreateNewProduct(payload InsertProductInput, userId string) (InsertProductResponse, int, error)
  GetAllProducts(page, size int) ([]ProductResponse, models.Paging, int, error)
  GetMerchantProducts(merchantId string, page, size int) ([]entities.Product, models.Paging, error)
  GetProductByID(id string) (ProductResponse, int, error)
  DeleteProductByID(id, userId string) (int, error)
}

type useCase struct {
  repository Repository
  usersUseCase users.UseCase
  jwtService services.JwtService
}

func (u *useCase) CreateNewProduct(payload InsertProductInput, userId string) (InsertProductResponse, int, error) {
  product := entities.Product{
    UserID: userId,
  }

  if err := copier.Copy(&product, &payload); err != nil {
    log.Println("products.usecase: CreateNewProduct.copier.Copy Err :", err)
    return InsertProductResponse{}, http.StatusInternalServerError, err
  }

  product, err := u.repository.Insert(product)
  if err != nil {
    log.Println("products.usecase: CreateNewProduct.repository.Insert Err :", err)
    return InsertProductResponse{}, http.StatusBadRequest, err
  }
  
  user, statusCode, err := u.usersUseCase.GetUserByID(product.UserID)
  if err != nil {
    log.Println("products.usecase: CreateNewProduct.usersUseCase.GetUserByID Err :", err)
    return InsertProductResponse{}, statusCode, err
  }

  userResponse := users.UserResponse{
    Name: user.Name,
    Username: user.Username,
    Email: user.Email,
  }

  productResponse := InsertProductResponse{
    ID: product.ID,
    Merchant: userResponse,
    ProductName: product.ProductName,
    Description: product.Description,
    Stock: product.Stock,
    Price: product.Price,
  } 

  return productResponse, http.StatusOK, nil
}

func (u *useCase) GetAllProducts(page, size int) ([]ProductResponse, models.Paging, int, error) {
  products, paging, err := u.repository.FindAll(page, size)
  if err != nil {
    log.Println("products.usecase: GetAllProducts.repository.FindAll Err :", err)
    return nil, models.Paging{}, http.StatusBadRequest, err
  }

  usersId := u.getUserIDFromProducts(products)
  usersMap := u.getUserByUserIDs(usersId)
  productResponses := u.getProductResponses(products, usersMap)
  
  return productResponses, paging, http.StatusOK, nil
}

func (u *useCase) GetMerchantProducts(merchantId string, page, size int) ([]entities.Product, models.Paging, error) {
  return u.repository.FindAllByUserID(merchantId, page, size)
}

func (u *useCase) GetProductByID(id string) (ProductResponse, int, error) {
  product, err := u.repository.FindByID(id)
  if err != nil {
    log.Println("products.usecase: GetProductByID.repository.FindByID Err :", err)
    return ProductResponse{}, http.StatusNotFound, err
  }

  user, statusCode, err := u.usersUseCase.GetUserByID(product.UserID)
  if err != nil {
    log.Println("products.usecase: GetProductByID.usersUseCase.GetUserByID Err :", err)
    return ProductResponse{}, statusCode, err
  }

  productResponse := u.getProductResponse(product, user)

  return productResponse, http.StatusOK, nil
}

func (u *useCase) DeleteProductByID(id, userId string) (int, error) {
  if err := u.repository.DeleteByID(id, userId); err != nil {
    return http.StatusNotFound, err
  }

  return http.StatusOK, nil
}

func NewUseCase(repository Repository, usersUseCase users.UseCase, jwtService services.JwtService) UseCase {
  return &useCase{repository, usersUseCase, jwtService}
}

func (u *useCase) getUserIDFromProducts(products []entities.Product) S {
  userId := S{}
  for _, product := range products {
    if !userId.contains(product.UserID) {
      userId = append(userId, product.UserID)
    }
  }

  return userId
}

func (u *useCase) getUserByUserIDs(usersId S) M {
  users := M{}
  for _, id := range usersId {
    user, _, err := u.usersUseCase.GetUserByID(id)
    if err != nil {
      log.Println("products.usecase: GetAllProducts.repository.FindByID Err :", err)
      return nil
    }

    users[id] = user
  }

  return users
}

func (u *useCase) getProductResponse(product entities.Product, user entities.User) ProductResponse {
  userResponse := users.UserResponse{
    Name: user.Name,
    Username: user.Username,
    Email: user.Email,
    PhoneNumber: user.PhoneNumber,
  }

  productResponse := ProductResponse{
    ID: product.ID,
    Merchant: userResponse,
    ProductName: product.ProductName,
    Description: product.Description,
    Stock: product.Stock,
    Price: product.Price,
    CreatedAt: product.CreatedAt,
    UpdatedAt: product.UpdatedAt,
  }

  return productResponse
}

func (u *useCase) getProductResponses(products []entities.Product, usersMap M) []ProductResponse {
  productResponses := []ProductResponse{}
  for _, product := range products {
    user := users.UserResponse{
      Name: usersMap[product.UserID].Name,
      Username: usersMap[product.UserID].Username,
      Email: usersMap[product.UserID].Email,
      PhoneNumber: usersMap[product.UserID].PhoneNumber,
    }
    
    productResponse := ProductResponse{
      ID: product.ID,
      Merchant: user,
      ProductName: product.ProductName,
      Description: product.Description,
      Stock: product.Stock,
      Price: product.Price,
      CreatedAt: product.CreatedAt,
      UpdatedAt: product.UpdatedAt,
    }
    
    productResponses = append(productResponses, productResponse)
  }

  return productResponses
}
