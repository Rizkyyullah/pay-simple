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
  GetAllProducts(page, size int) ([]GetAllProductsResponse, models.Paging, int, error)
}

type useCase struct {
  repository Repository
  usersRepository users.Repository
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
  
  user, err := u.usersRepository.FindByID(product.UserID)
  if err != nil {
    log.Println("products.usecase: CreateNewProduct.usersRepository.FindByID Err :", err)
    return InsertProductResponse{}, http.StatusBadRequest, err
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

func (u *useCase) GetAllProducts(page, size int) ([]GetAllProductsResponse, models.Paging, int, error) {
  products, paging, err := u.repository.FindAll(page, size)
  if err != nil {
    log.Println("products.usecase: GetAllProducts.repository.FindAll Err :", err)
    return nil, models.Paging{}, http.StatusBadRequest, err
  }

  userId := S{}
  for _, product := range products {
    if !userId.contains(product.UserID) {
      userId = append(userId, product.UserID)
    }
  }

  usersMap := map[string]entities.User{}
  for _, id := range userId {
    user, err := u.usersRepository.FindByID(id)
    if err != nil {
      log.Println("products.usecase: GetAllProducts.repository.FindByID Err :", err)
      return nil, models.Paging{}, http.StatusBadRequest, err
    }

    usersMap[id] = user
  }

  productResponses := []GetAllProductsResponse{}
  for _, product := range products {
    user := users.UserResponse{
      Name: usersMap[product.UserID].Name,
      Username: usersMap[product.UserID].Username,
      Email: usersMap[product.UserID].Email,
      PhoneNumber: usersMap[product.UserID].PhoneNumber,
    }
    
    productResponse := GetAllProductsResponse{
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
  
  return productResponses, paging, http.StatusOK, nil
}

func NewUseCase(repository Repository, usersRepository users.Repository, jwtService services.JwtService) UseCase {
  return &useCase{repository, usersRepository, jwtService}
}
