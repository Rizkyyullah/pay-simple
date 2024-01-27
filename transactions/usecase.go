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
  CreateTransaction(userId string, payload CreateTransactionInput) error
  Topup(payload TopupInput, userId string) (users.GetBalanceResponse, error)
  Transfer(payload TransferInput, userId string) (users.GetBalanceResponse, error)
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

func (u *useCase) CreateTransaction(userId string, payload CreateTransactionInput) error {
  paidStatus := true
  cashflow := "MONEY_OUT"
  transactionType := "DEBIT"

  products, err := u.getProducts(payload)
  if err != nil {
    return err
  }

  merchants, err := u.getMerchants(payload)
  if err != nil {
    return err
  }

  quantities := u.getQuantities(payload)
  balance, _, err := u.usersUC.GetBalance(userId)
  if err != nil {
    return err
  }

  dto := TransactionDTO{
    UserID: userId,
    TransactionType: transactionType,
    PaidStatus: paidStatus,
    Cashflow: cashflow,
    Products: products,
    Merchants: merchants,
    Quantities: quantities,
    Balance: balance.Balance,
  }

  return u.repository.Insert(dto)
}

func (u *useCase) Topup(payload TopupInput, userId string) (users.GetBalanceResponse, error) {
  balance, _, err := u.usersUC.GetBalance(userId)
  if err != nil {
    return users.GetBalanceResponse{}, err
  }
  
  dto := TransactionDTO{
    UserID: userId,
    TransactionType: "CREDIT",
    PaidStatus: false,
    Cashflow: "MONEY_IN",
    Balance: balance.Balance + payload.Amount,
    Transfer: false,
  }

  return u.repository.UpdateBalance(dto, "")
}

func (u *useCase) Transfer(payload TransferInput, userId string) (users.GetBalanceResponse, error) {
  balance, _, err := u.usersUC.GetBalance(userId)
  if err != nil {
    return users.GetBalanceResponse{}, err
  }
  
  dto := TransactionDTO{
    UserID: userId,
    TransactionType: "DEBIT",
    PaidStatus: false,
    Cashflow: "MONEY_OUT",
    Balance: balance.Balance - payload.Amount,
    Amount: payload.Amount,
    Transfer: true,
  }

  return u.repository.UpdateBalance(dto, payload.ToUserId)
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
