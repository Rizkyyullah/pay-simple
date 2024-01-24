package users

import (
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "log"
  "net/http"
)

type UseCase interface {
  GetBalance(id string) (GetBalanceResponse, int, error)
}

type useCase struct {
  repository Repository
  jwtService services.JwtService
}

func (u *useCase) GetBalance(id string) (GetBalanceResponse, int, error) {
  user, err := u.repository.FindByID(id)
  if err != nil {
    log.Println("users.UseCase: GetBalance.FindByID Err :", err)
    return GetBalanceResponse{}, http.StatusBadRequest, err
  }

  return GetBalanceResponse{user.Balance}, http.StatusOK, nil
}

func NewUseCase(repository Repository, jwtService services.JwtService) UseCase {
  return &useCase{repository, jwtService}
}
