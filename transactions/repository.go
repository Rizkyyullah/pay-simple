package transactions

import (
  "context"
  "errors"
  "fmt"
  "log"
  "math"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  
  "github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Repository interface {
  FindAllByUserID(userId string, page, size int) ([]entities.Transaction, models.Paging, error)
  FindByID(id, userId string) (entities.Transaction, error)
}

type repository struct {
  conn *pgx.Conn
}

func (r *repository) FindAllByUserID(userId string, page, size int) ([]entities.Transaction, models.Paging, error) {
  var transactions []entities.Transaction
  var offset = (page - 1) * size

  rows, err := r.conn.Query(ctx, configs.SelectAllTransactionByUserID, userId, size, offset)
  if err != nil {
    log.Println("transactions.repository: FindAllByUserID.Query Err :", err)
    return nil, models.Paging{}, err
  }
  defer rows.Close()

  for rows.Next() {
    var transaction entities.Transaction
    if err := rows.Scan(
      &transaction.ID,
      &transaction.UserID,
      &transaction.TransactionDate,
      &transaction.TransactionType,
      &transaction.PaidStatus,
      &transaction.Cashflow,
      &transaction.CreatedAt,
      &transaction.UpdatedAt,
    ); err != nil {
      if errors.Is(err, pgx.ErrNoRows) {
        return nil, models.Paging{}, fmt.Errorf("You don't have a transaction history")
      } else {
        return nil, models.Paging{}, err
      }
    }
    
    transactions = append(transactions, transaction)
  }

  var totalRows = 0
  if err := r.conn.QueryRow(ctx, "SELECT COUNT(*) FROM transactions WHERE user_id = $1", userId).Scan(&totalRows); err != nil {
    log.Println("transactions.Repository: FindAllByUserID.QueryRow.Scan Err :", err)
    return nil, models.Paging{}, err
  }

  paging := models.Paging{
    Page: page,
    RowsPerPage: size,
    TotalRows: totalRows,
    TotalPages: int(math.Ceil(float64(totalRows) / float64(size))),
  }

  return transactions, paging, nil
}

func (r *repository) FindByID(id, userId string) (entities.Transaction, error) {
  var transaction entities.Transaction

  if err := r.conn.QueryRow(ctx, configs.SelectTransactionByID, id, userId).Scan(
    &transaction.ID,
    &transaction.TransactionDate,
    &transaction.TransactionType,
    &transaction.PaidStatus,
    &transaction.Cashflow,
    &transaction.CreatedAt,
  ); err != nil {
    log.Println("transactions.repository: FindByID.QueryRow.Scan Err :", err)
    return entities.Transaction{}, err
  }

  return transaction, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}