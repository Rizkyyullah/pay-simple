package users

import (
  "log"
  "github.com/Rizkyyullah/pay-simple/entities"
  "strings"

  "github.com/jinzhu/copier"
  "golang.org/x/crypto/bcrypt"
)

type UseCase interface {
  Register(payload RegisterInput, path string) (entities.User, error)
}

type useCase struct {
  repository Repository
}

func (u *useCase) Register(payload RegisterInput, path string) (entities.User, error) {
  hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
  if err != nil {
    log.Println("users.UseCase: GenerateFromPassword Err :", err)
    return entities.User{}, err
  }
  
  var user entities.User
  if err := copier.Copy(&user, &payload); err != nil {
    log.Println("users.UseCase: copier.Copy Err :", err)
    return entities.User{}, err
  }
  
  user.Password = string(hashPassword)
  
  if strings.Contains(path, "merchants") {
    user.Role = "MERCHANT"
  } else {
    user.Role = "CUSTOMER"
  }
  
  return u.repository.Insert(user)
}

func NewUseCase(repository Repository) UseCase {
  return &useCase{repository}
}
