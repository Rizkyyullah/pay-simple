package entities

import "time"

type Product struct {
  ID          string    `json:"id"`
  UserID      string    `json:"userId,omitempty"`
  ProductName string    `json:"productName,omitempty"`
  Description string    `json:"description,omitempty"`
  Stock       int    `json:"stock,omitempty"`
  Price       int    `json:"price,omitempty"`
  CreatedAt   time.Time `json:"createdAt,omitempty"`
  UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
