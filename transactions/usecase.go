package transactions

import (
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "github.com/Rizkyyullah/pay-simple/transactions_detail"
)

type UseCase interface {
  GetTransactionsHistory(userId string, page, size int) ([]TransactionResponse, models.Paging, error)
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

  transactionsResponse := u.getTransactionsResponse(transactions)

  return transactionsResponse, paging, nil
}

func NewUseCase(repository Repository, transactionsDetailUC transactions_detail.UseCase) UseCase {
  return &useCase{repository, transactionsDetailUC}
}

func (u *useCase) getTransactionsResponse(transactions []entities.Transaction) []TransactionResponse {
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
