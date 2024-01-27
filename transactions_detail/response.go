package transactions_detail

import "github.com/Rizkyyullah/pay-simple/products"

type TransactionDetailResponse struct {
  ID            string                      `json:"id"`
  TransactionID string                      `json:"transactionId,omitempty"`
  Products      []products.ProductResponse  `json:"products,omitempty"`
  Quantity      int                         `json:"quantity,omitempty"`
  TotalPrice    int                         `json:"totalPrice,omitempty"`
  CreatedAt     string                      `json:"createdAt,omitempty"`
  UpdatedAt     string                      `json:"updatedAt,omitempty"`
}