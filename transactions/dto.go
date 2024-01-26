package transactions

import (
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/products"
)

type TransactionDTO struct {
  UserID          string
  TransactionType string
  PaidStatus      bool
  Cashflow        string
  Products        []products.ProductResponse
  Merchants       map[string]entities.User
  Quantities      map[string]int
  Balance         int
}