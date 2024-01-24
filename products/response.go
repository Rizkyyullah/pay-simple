package products

import (
  "github.com/Rizkyyullah/pay-simple/users"
  "time"
)

type InsertProductResponse struct {
  ID          string              `json:"id"`
  Merchant    users.UserResponse  `json:"merchant,omitempty"`
  ProductName string              `json:"productName,omitempty"`
  Description string              `json:"description,omitempty"`
  Stock       int                 `json:"stock,omitempty"`
  Price       int                 `json:"price,omitempty"`
}

type ProductResponse struct {
  ID          string              `json:"id"`
  Merchant    users.UserResponse  `json:"merchant,omitempty"`
  ProductName string              `json:"productName,omitempty"`
  Description string              `json:"description,omitempty"`
  Stock       int                 `json:"stock,omitempty"`
  Price       int                 `json:"price,omitempty"`
  CreatedAt   time.Time           `json:"createdAt,omitempty"`
  UpdatedAt   time.Time           `json:"updatedAt,omitempty"`
}