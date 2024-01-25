package transactions

import td "github.com/Rizkyyullah/pay-simple/transactions_detail"

type TransactionResponse struct {
  ID                  string                          `json:"id"`
  TransactionDetails  []td.TransactionDetailResponse  `json:"transactionDetails,omitempty"`
  TransactionDate     string                          `json:"transactionDate,omitempty"`
  TransactionType     string                          `json:"transactionType,omitempty"`
  PaidStatus          bool                            `json:"paidStatus,omitempty"`
  Cashflow            string                          `json:"cashflow,omitempty"`
  CreatedAt           string                          `json:"createdAt,omitempty"`
  UpdatedAt           string                          `json:"updatedAt,omitempty"`
}