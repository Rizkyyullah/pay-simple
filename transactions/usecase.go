package transactions

import (
  "errors"
  "fmt"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/products"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "github.com/Rizkyyullah/pay-simple/transactions_detail"
  "github.com/Rizkyyullah/pay-simple/users"
  "net/http"

  "github.com/jackc/pgx/v5"
)

type UseCase interface {
  GetTransactionsHistory(userId string, page, size int) ([]TransactionResponse, models.Paging, error)
  GetTransactionHistoryByID(id, userId string) (TransactionResponse, int, error)
  CreateTransaction(userId string, payload CreateTransactionInput) (TransactionResponse, error)
}

type useCase struct {
  repository Repository
  transactionsDetailUC transactions_detail.UseCase
  productsUC products.UseCase
  usersUC users.UseCase
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

func (u *useCase) CreateTransaction(userId string, payload CreateTransactionInput) (TransactionResponse, error) {
  paidStatus := true
  cashflow := "MONEY_OUT"
  transactionType := "DEBIT"

  products, err := u.getProducts(payload)
  if err != nil {
    return TransactionResponse{}, err
  }

  merchants, err := u.getMerchants(payload)
  if err != nil {
    return TransactionResponse{}, err
  }

  quantities := u.getQuantities(payload)
  balance, _, err := u.usersUC.GetBalance(userId)
  if err != nil {
    return TransactionResponse{}, err
  }

  transaction, err := u.repository.Insert(userId, transactionType, paidStatus, cashflow, products, merchants, quantities, balance.Balance)
  if err != nil {
    return TransactionResponse{}, err
  }

  trResponse := TransactionResponse{
    ID: transaction.ID,
    TransactionDetails: nil,
    TransactionDate: transaction.TransactionDate.Format("Monday, 02 January 2006"),
    TransactionType: transaction.TransactionType,
    PaidStatus: transaction.PaidStatus,
    Cashflow: transaction.Cashflow,
    CreatedAt: transaction.CreatedAt.Format("Monday, 02 January 2006 15:04:05"),
  }

  return trResponse, nil
}

func NewUseCase(repository Repository, transactionsDetailUC transactions_detail.UseCase, productsUC products.UseCase, usersUC users.UseCase) UseCase {
  return &useCase{repository, transactionsDetailUC, productsUC, usersUC}
}

func (u *useCase) getTransactionsResponse(transactions ...entities.Transaction) []TransactionResponse {
  var transactionsResponse []TransactionResponse
  for _, val := range transactions {
    transactionsDetail, _ := u.transactionsDetailUC.GetTransactionsDetail(val.ID)
    transactionResponse := TransactionResponse{
      ID: val.ID,
      TransactionDetails: transactionsDetail,
      TransactionDate: val.TransactionDate.Format("Monday, 02 January 2006"),
      TransactionType: val.TransactionType,
      PaidStatus: val.PaidStatus,
      Cashflow: val.Cashflow,
      CreatedAt: val.CreatedAt.Format("Monday, 02 January 2006 15:04:05 WIB"),
    }

    transactionsResponse = append(transactionsResponse, transactionResponse)
  }

  return transactionsResponse
}

func (u *useCase) getProducts(payload CreateTransactionInput) ([]products.ProductResponse, error) {
  var products []products.ProductResponse
  for _, val := range payload.Products {
    product, _, err := u.productsUC.GetProductByID(val.ID)
    if err != nil {
      return nil, err
    }

    products = append(products, product)
  }

  return products, nil
}

func (u *useCase) getMerchants(payload CreateTransactionInput) (map[string]entities.User, error) {
  merchants := map[string]entities.User{}
  for _, val := range payload.Products {
    merchant, _, err := u.usersUC.GetUserByID(val.MerchantID)
    if err != nil {
      return nil, err
    }

    merchants[val.ID] = merchant
  }

  return merchants, nil
}

func (u *useCase) getQuantities(payload CreateTransactionInput) map[string]int {
  quantities := map[string]int{}
  for _, val := range payload.Products {
    quantities[val.ID] = val.Quantity
  }

  return quantities
}
