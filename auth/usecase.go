package auth

import (
  "fmt"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/users"
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "log"
  "net/http"
  "strings"

  "github.com/jinzhu/copier"
  "golang.org/x/crypto/bcrypt"
)

type UseCase interface {
  Register(payload RegisterInput, path string) (RegisterResponse, int, error)
  Login(payload LoginInput, path string) (string, int, error)
}

type useCase struct {
  usersRepository users.Repository
  jwtService services.JwtService
}

func (u *useCase) Register(payload RegisterInput, path string) (RegisterResponse, int, error) {
  hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
  if err != nil {
    log.Println("auth.UseCase: GenerateFromPassword Err :", err)
    return RegisterResponse{}, http.StatusInternalServerError, err
  }
  
  var user entities.User
  if err := copier.Copy(&user, &payload); err != nil {
    log.Println("auth.UseCase: copier.Copy Err :", err)
    return RegisterResponse{}, http.StatusInternalServerError, err
  }
  
  user.Password = string(hashPassword)
  
  if strings.HasSuffix(path, "merchants") {
    user.Role = "MERCHANT"
  } else {
    user.Role = "CUSTOMER"
  }
  
  user, err = u.usersRepository.Insert(user)
  if err != nil {
    log.Println("auth.UseCase: usersRepository.Insert Err :", err)
    return RegisterResponse{}, http.StatusInternalServerError, err
  }
  
  var userResponse RegisterResponse
  if err := copier.Copy(&userResponse, &user); err != nil {
    return RegisterResponse{}, http.StatusInternalServerError, err
  }
  
  return userResponse, http.StatusCreated, nil
}

func (u *useCase) Login(payload LoginInput, path string) (string, int, error) {
  user, err := u.usersRepository.FindByEmail(payload.Email)
  if err != nil {
    log.Println("auth.UseCase: FindByEmail Err :", err)
    return "", http.StatusNotFound, err
  }
  
  expectedRole := "CUSTOMER"
  errorMsg := "customer"
  if strings.HasSuffix(path, "merchants") {
    expectedRole = "MERCHANT"
    errorMsg = "merchant"
  }

  if user.Role != expectedRole {
    return "", http.StatusBadRequest, fmt.Errorf("User with email '%s' is not a %s", payload.Email, errorMsg)
  }
  
  if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)) != nil {
    log.Println("auth.UseCase: CompareHashAndPassword Err")
    return "", http.StatusBadRequest, fmt.Errorf("The password you entered is wrong")
  }
  
  token, err := u.jwtService.CreateToken(user)
  if err != nil {
    log.Println("auth.UseCase: CreateToken Err :", err)
    return "", http.StatusInternalServerError, err
  }

  return token, http.StatusOK, nil
}

func NewUseCase(usersRepository users.Repository, jwtService services.JwtService) UseCase {
  return &useCase{usersRepository, jwtService}
}
