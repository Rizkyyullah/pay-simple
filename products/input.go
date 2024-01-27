package products

type InsertProductInput struct {
  ProductName string  `json:"productName" validate:"required,min=5"`
  Description string  `json:"description"`
  Stock       int     `json:"stock" validate:"required,number,gt=0"`
  Price       int     `json:"price" validate:"required,number,gt=0"`
}