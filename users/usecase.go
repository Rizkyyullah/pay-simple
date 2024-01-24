package users

import (
  "errors"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/services"
  "log"
  "net/http"

  "github.com/jackc/pgx/v5"
)

type UseCase interface {
  GetBalance(id string) (GetBalanceResponse, int, error)
  GetUserByID(id string) (entities.User, int, error)
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

func (u *useCase) GetUserByID(id string) (entities.User, int, error) {
  user, err := u.repository.FindByID(id)
  if err != nil {
    log.Println("users.UseCase: GetUserByID.FindByID Err :", err)
    if errors.Is(err, pgx.ErrNoRows) {
      return entities.User{}, http.StatusNotFound, err
    } else {
      return entities.User{}, http.StatusInternalServerError, err
    }
  }

  return user, http.StatusOK, nil
}

func NewUseCase(repository Repository, jwtService services.JwtService) UseCase {
  return &useCase{repository, jwtService}
}
