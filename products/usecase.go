package products

import (
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "github.com/Rizkyyullah/pay-simple/users"
  "log"
  "net/http"
  
  "github.com/jinzhu/copier"
)

type UseCase interface {
  CreateNewProduct(payload InsertProductInput, userId string) (InsertProductResponse, int, error)
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

func NewUseCase(repository Repository, usersRepository users.Repository, jwtService services.JwtService) UseCase {
  return &useCase{repository, usersRepository, jwtService}
}
