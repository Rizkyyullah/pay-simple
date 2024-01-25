package entities

import "time"

type Transaction struct {
  ID              string    `json: "id"`
  UserID          string    `json: "userId"`
  TransactionDate time.Time `json: "transactionDate"`
  TransactionType string    `json: "transactionType"`
  PaidStatus      bool      `json: "paidStatus"`
  Cashflow        string    `json: "cashFlow"`
  CreatedAt       time.Time `json: "createdAt"`
  UpdatedAt       time.Time `json: "updatedAt"`
}