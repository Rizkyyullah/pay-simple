package transactions_detail

import (
  "fmt"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/products"
  "github.com/Rizkyyullah/pay-simple/shared/common"
)

type UseCase interface {
  GetTransactionsDetail(transactionId string) ([]TransactionDetailResponse, error)
}

type useCase struct {
  repository Repository
  productsUseCase products.UseCase
}

func (u *useCase) GetTransactionsDetail(transactionId string) ([]TransactionDetailResponse, error) {
  fmt.Println(transactionId)
  transactionsDetail, err := u.repository.FindAllByTransactionIDAndProductID(transactionId)
  if err != nil {
    return nil, err
  }
  fmt.Println(transactionsDetail)

  transactionsDetailResponse := u.getTransactionsDetailResponse(transactionsDetail)
  return transactionsDetailResponse, nil
}

func NewUseCase(repository Repository, productsUseCase products.UseCase) UseCase {
  return &useCase{repository, productsUseCase}
}

func (u *useCase) getTransactionsDetailResponse(transactionsDetail []entities.TransactionsDetail) []TransactionDetailResponse {
  var transactionsDetailResponse []TransactionDetailResponse
  for _, val := range transactionsDetail {
    tdResponse := TransactionDetailResponse {
      ID: val.ID,
      Quantity: val.Quantity,
      TotalPrice: val.TotalPrice,
      CreatedAt: val.CreatedAt.In(common.GetTimezone()).Format("Monday, 02 January 2006 15:04:05 MST"),
    }

    transactionsDetailResponse = append(transactionsDetailResponse, tdResponse)
  }

  return transactionsDetailResponse
}
