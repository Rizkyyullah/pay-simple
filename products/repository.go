package products

import (
  "context"
  "log"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  
  "github.com/jackc/pgx/v5"
)

type Repository interface {
  Insert(payload entities.Product) (entities.Product, error)
}

type repository struct {
  conn *pgx.Conn
}

func (r *repository) Insert(payload entities.Product) (entities.Product, error) {
  var product entities.Product
  
  id, err := common.UniqueID(r.conn, "PRD", "products_id_seq")
  if err != nil {
    log.Println("users.Repository: UniqueID Err :", err)
    return entities.Product{}, err
  }

  if err := r.conn.QueryRow(context.Background(), configs.InsertProduct,
    id,
    payload.UserID,
    payload.ProductName,
    payload.Description,
    payload.Stock,
    payload.Price,
  ).Scan(
    &product.ID,
    &product.UserID,
    &product.ProductName,
    &product.Description,
    &product.Stock,
    &product.Price,
  ); err != nil {
    log.Println("products.repository: Insert.QueryRow.Scan Err :", err)
    return entities.Product{}, err
  }

  return product, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}