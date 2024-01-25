package entities

import "time"

type TransactionsDetail struct {
  ID            string    `json:"id"`
  TransactionID string    `json:"transactionId"`
  ProductID     string    `json:"productId"`
  Quantity      int       `json:"quantity"`
  TotalPrice    int       `json:"totalPrice"`
  CreatedAt     time.Time `json:"createdAt"`
  UpdatedAt     time.Time `json:"updatedAt"`
}