package products

import (
  "context"
  "errors"
  "fmt"
  "log"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "math"
  
  "github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Repository interface {
  Insert(payload entities.Product) (entities.Product, error)
  FindAll(page, size int) ([]entities.Product, models.Paging, error)
  FindAllByUserID(userId string, page, size int) ([]entities.Product, models.Paging, error)
  FindByID(id string) (entities.Product, error)
  DeleteByID(id, userId string) error
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

  if err := r.conn.QueryRow(ctx, configs.InsertProduct,
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
  
  rows, err := r.conn.Query(ctx, configs.SelectAllProducts, size, offset)
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
  if err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM products").Scan(&totalRows); err != nil {
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

func (r *repository) FindAllByUserID(userId string, page, size int) ([]entities.Product, models.Paging, error) {
  var products []entities.Product
  var offset = (page - 1) * size
  
  rows, err := r.conn.Query(ctx, configs.SelectAllProductsByUserID, userId, size, offset)
  if err != nil {
    log.Println("products.Repository: FindByID.Query Err :", err)
    return nil, models.Paging{}, err
  }
  defer rows.Close()

  for rows.Next() {
    var product entities.Product
    if err := rows.Scan(
      &product.ID,
      &product.ProductName,
      &product.Description,
      &product.Stock,
      &product.Price,
      &product.CreatedAt,
      &product.UpdatedAt,
    ); err != nil {
      log.Println("products.Repository: FindByID.Scan Err :", err)
      return nil, models.Paging{}, err
    }

    products = append(products, product)
  }

  var totalRows = 0
  if err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM products WHERE user_id = $1", userId).Scan(&totalRows); err != nil {
    log.Println("products.Repository: FindByID.QueryRow Err :", err)
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

func (r *repository) FindByID(id string) (entities.Product, error) {
  var product entities.Product

  if err := r.conn.QueryRow(ctx, configs.SelectProductByID, id).Scan(
    &product.ID,
    &product.UserID,
    &product.ProductName,
    &product.Description,
    &product.Stock,
    &product.Price,
    &product.CreatedAt,
    &product.UpdatedAt,
  ); err != nil {
    log.Println("products.repository: FindByID.QueryRow.Scan Err :", err)
    if errors.Is(err, pgx.ErrNoRows) {
      return entities.Product{}, fmt.Errorf("Product with id '%s' not found", id)
    } else {
      return entities.Product{}, err
    }
  }

  return product, nil
}

func (r *repository) DeleteByID(id, userId string) error {
  if err := r.conn.QueryRow(ctx, configs.DeleteProductByID, id, userId).Scan(&id); err != nil {
    if errors.Is(err, pgx.ErrNoRows) {
      return fmt.Errorf("Product with id '%s' not found", id)
    } else {
      return err
    }
  }

  return nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}