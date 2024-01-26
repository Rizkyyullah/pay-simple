package transactions

type ProductsInput struct {
  ID          string  `json:"id" validate:"required"`
  MerchantID  string  `json:"merchantId" validate:"required"`
  Quantity    int     `json:"quantity" validate:"required,number,min=1"`
}

type CreateTransactionInput struct {
  Products []ProductsInput `json:"products"`
}

type TopupInput struct {
  Amount  int     `json:"amount" validate:"required,number,min=5000"`
}