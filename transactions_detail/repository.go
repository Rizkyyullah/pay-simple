package transactions_detail

import (
  "context"
  "log"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  
  "github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Repository interface {
  FindAllByTransactionIDAndProductID(transactionId string) ([]entities.TransactionsDetail, error)
}

type repository struct {
  conn *pgx.Conn
}

func (r *repository) FindAllByTransactionIDAndProductID(transactionId string) ([]entities.TransactionsDetail, error) {
  var transactionsDetail []entities.TransactionsDetail

  rows, err := r.conn.Query(ctx, configs.SelectAllTransactionsDetail, transactionId)
  if err != nil {
    log.Println("transactions_detail.repository: FindAllByTransactionIDAndProductID.Query Err :", err)
    return nil, err
  }
  defer rows.Close()

  for rows.Next() {
    var transactionDetail entities.TransactionsDetail
    if err := rows.Scan(
      &transactionDetail.ID,
      &transactionDetail.ProductID,
      &transactionDetail.Quantity,
      &transactionDetail.TotalPrice,
      &transactionDetail.CreatedAt,
    ); err != nil {
      return nil, err
    }
    
    transactionsDetail = append(transactionsDetail, transactionDetail)
  }

  var totalRows = 0
  if err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM transaction_details WHERE transaction_id = $1", transactionId).Scan(&totalRows); err != nil {
    log.Println("transactions_detail.Repository: FindAllByTransactionIDAndProductID.QueryRow.Scan Err :", err)
    return nil, err
  }

  return transactionsDetail, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}