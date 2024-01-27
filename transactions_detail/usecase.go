package transactions_detail

import (
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/products"
)

type UseCase interface {
  GetTransactionsDetail(transactionId string) ([]TransactionDetailResponse, error)
  CreateTransactionDetail(id, transactionId, productId string, quantity, totalPrice int) (TransactionDetailResponse, error)
}

type useCase struct {
  repository Repository
  productsUseCase products.UseCase
}

func (u *useCase) GetTransactionsDetail(transactionId string) ([]TransactionDetailResponse, error) {
  transactionsDetail, err := u.repository.FindAllByTransactionIDAndProductID(transactionId)
  if err != nil {
    return nil, err
  }

  transactionsDetailResponse := u.getTransactionsDetailResponse(transactionsDetail)
  return transactionsDetailResponse, nil
}

func (u *useCase) CreateTransactionDetail(id, transactionId, productId string, quantity, totalPrice int) (TransactionDetailResponse, error) {
  transactionDetail, err := u.repository.Insert(id, transactionId, productId, quantity, totalPrice)
  if err != nil {
    return TransactionDetailResponse{}, err
  }

  tdResponse := TransactionDetailResponse{
    ID: transactionDetail.ID,
    TransactionID: transactionDetail.TransactionID,
    Quantity: transactionDetail.Quantity,
    TotalPrice: transactionDetail.TotalPrice,
    CreatedAt: transactionDetail.CreatedAt.Format("Monday, 02 January 2006 15:04:05 WIB"),
  }

  return tdResponse, nil
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
      CreatedAt: val.CreatedAt.Format("Monday, 02 January 2006 15:04:05 WIB"),
    }

    transactionsDetailResponse = append(transactionsDetailResponse, tdResponse)
  }

  return transactionsDetailResponse
}
