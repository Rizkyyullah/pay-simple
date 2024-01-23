package customers

import (
  "github.com/Rizkyyullah/pay-simple/auth"
  "github.com/Rizkyyullah/pay-simple/users"
  "log"
  "net/http"
)

type UseCase interface {
  GetBalance(id string) (GetBalanceResponse, int, error)
}

type useCase struct {
  usersRepository users.Repository
  jwtService auth.JwtService
}

func (u *useCase) GetBalance(id string) (GetBalanceResponse, int, error) {
  user, err := u.usersRepository.FindByID(id)
  if err != nil {
    log.Println("customers.UseCase: GetBalance.FindByIDAndRole Err :", err)
    return GetBalanceResponse{}, http.StatusBadRequest, err
  }

  return GetBalanceResponse{user.Balance}, http.StatusOK, nil
}

func NewUseCase(usersRepository users.Repository, jwtService auth.JwtService) UseCase {
  return &useCase{usersRepository, jwtService}
}
