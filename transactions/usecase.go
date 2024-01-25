package transactions

import (
  "errors"
  "fmt"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "github.com/Rizkyyullah/pay-simple/transactions_detail"
  "net/http"

  "github.com/jackc/pgx/v5"
)

type UseCase interface {
  GetTransactionsHistory(userId string, page, size int) ([]TransactionResponse, models.Paging, error)
  GetTransactionHistoryByID(id, userId string) (TransactionResponse, int, error)
}

type useCase struct {
  repository Repository
  transactionsDetailUC transactions_detail.UseCase
}

func (u *useCase) GetTransactionsHistory(userId string, page, size int) ([]TransactionResponse, models.Paging, error) {
  transactions, paging, err := u.repository.FindAllByUserID(userId, page, size)
  if err != nil {
    return nil, models.Paging{}, err
  }

  transactionsResponse := u.getTransactionsResponse(transactions...)
  return transactionsResponse, paging, nil
}

func (u *useCase) GetTransactionHistoryByID(id, userId string) (TransactionResponse, int, error) {
  transaction, err := u.repository.FindByID(id, userId)
  if err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
      return TransactionResponse{}, http.StatusNotFound, fmt.Errorf("Transaction with id '%s' not found", id)
    } else {
      return TransactionResponse{}, http.StatusBadRequest, err
    }
  }

  transactionResponse := u.getTransactionsResponse(transaction)[0]
  return transactionResponse, http.StatusOK, nil
}

func NewUseCase(repository Repository, transactionsDetailUC transactions_detail.UseCase) UseCase {
  return &useCase{repository, transactionsDetailUC}
}

// func (u *useCase) getTransactionResponse(transaction entities.Transaction) []TransactionResponse {
//   transactionsDetail, _ := u.transactionsDetailUC.GetTransactionsDetail(transaction.ID)
//   transactionResponse := TransactionResponse{
//     ID: val.ID,
//     TransactionDetails: transactionsDetail,
//     TransactionDate: val.TransactionDate.In(common.GetTimezone()).Format("Monday, 02 January 2006"),
//     TransactionType: val.TransactionType,
//     PaidStatus: val.PaidStatus,
//     Cashflow: val.Cashflow,
//     CreatedAt: val.CreatedAt.In(common.GetTimezone()).Format("Monday, 02 January 2006 15:04:05 MST"),
//   }

//   transactionsResponse = append(transactionsResponse, transactionResponse)

//   return transactionsResponse
// }

func (u *useCase) getTransactionsResponse(transactions ...entities.Transaction) []TransactionResponse {
  var transactionsResponse []TransactionResponse
  for _, val := range transactions {
    transactionsDetail, _ := u.transactionsDetailUC.GetTransactionsDetail(val.ID)
    transactionResponse := TransactionResponse{
      ID: val.ID,
      TransactionDetails: transactionsDetail,
      TransactionDate: val.TransactionDate.In(common.GetTimezone()).Format("Monday, 02 January 2006"),
      TransactionType: val.TransactionType,
      PaidStatus: val.PaidStatus,
      Cashflow: val.Cashflow,
      CreatedAt: val.CreatedAt.In(common.GetTimezone()).Format("Monday, 02 January 2006 15:04:05 MST"),
    }

    transactionsResponse = append(transactionsResponse, transactionResponse)
  }

  return transactionsResponse
}
