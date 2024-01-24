package products

import (
  "context"
  "log"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "math"
  
  "github.com/jackc/pgx/v5"
)

type Repository interface {
  Insert(payload entities.Product) (entities.Product, error)
  FindAll(page, size int) ([]entities.Product, models.Paging, error)
}

type repository struct {
  conn *pgx.Conn
}

func (r *repository) Insert(payload entities.Product) (entities.Product, error) {
  var product entities.Product
  
  id, err := common.UniqueID(r.conn, "PRD", "products_id_seq")
  if err != nil {
    log.Println("products.Repository: UniqueID Err :", err)
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

func (r *repository) FindAll(page, size int) ([]entities.Product, models.Paging, error) {
  var products []entities.Product
  var offset = (page - 1) * size
  
  rows, err := r.conn.Query(context.Background(),configs.SelectAllProducts, size, offset)
  if err != nil {
    log.Println("products.Repository: FindAll.Query Err :", err)
    return nil, models.Paging{}, err
  }
  defer rows.Close()

  for rows.Next() {
    var product entities.Product
    if err := rows.Scan(
      &product.ID,
      &product.UserID,
      &product.ProductName,
      &product.Description,
      &product.Stock,
      &product.Price,
      &product.CreatedAt,
      &product.UpdatedAt,
    ); err != nil {
      log.Println("products.Repository: FindAll.Scan Err :", err)
      return nil, models.Paging{}, err
    }

    products = append(products, product)
  }

  var totalRows = 0
  if err := r.conn.QueryRow(context.Background(), configs.SelectTotalRows).Scan(&totalRows); err != nil {
    log.Println("products.Repository: FindAll.QueryRow Err :", err)
    return nil, models.Paging{}, err
  }

  paging := models.Paging{
    Page: page,
    RowsPerPage: size,
    TotalRows: totalRows,
    TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
  }

  return products, paging, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}