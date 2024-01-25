package transactions

import (
  "context"
  "errors"
  "fmt"
  "log"
  "math"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/products"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  
  "github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Repository interface {
  FindAllByUserID(userId string, page, size int) ([]entities.Transaction, models.Paging, error)
  FindByID(id, userId string) (entities.Transaction, error)
  Insert(userId, transactionType string, paidStatus bool, cashflow string, products []products.ProductResponse, merchants map[string]entities.User, quantities map[string]int, balance int) (entities.Transaction, error)
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

func (r *repository) Insert(userId, transactionType string, paidStatus bool, cashflow string, products []products.ProductResponse, merchants map[string]entities.User, quantities map[string]int, balance int) (entities.Transaction, error) {
  var transaction entities.Transaction

  tx, err := r.conn.Begin(ctx)
  if err != nil {
    return entities.Transaction{}, err
  }
  defer tx.Rollback(ctx)
  
  trxId, _ := common.UniqueID(tx.Conn(), "TRX", "transactions_id_seq")
  if err := tx.QueryRow(ctx, configs.InsertTransaction, trxId, userId, transactionType, paidStatus, cashflow).Scan(
    &transaction.ID,
    &transaction.UserID,
    &transaction.TransactionDate,
    &transaction.TransactionType,
    &transaction.PaidStatus,
    &transaction.Cashflow,
    &transaction.CreatedAt,
  ); err != nil {
    log.Println("transactions.repository: Insert.QueryRow.Scan Err :", err)
    return entities.Transaction{}, err
  }

  currentBalance := balance
  for _, product := range products {
    productId := product.ID
    quantity := quantities[productId]
    totalPrice := product.Price * quantity
    merchantId := merchants[product.ID].ID
    currentBalance -= totalPrice
    if currentBalance < 0 {
      return entities.Transaction{}, fmt.Errorf("Your balance is insufficient")
    }

    tdId, _ := common.UniqueID(tx.Conn(), "TRXD", "transaction_details_id_seq")
    if _, err := tx.Exec(ctx, configs.InsertTransactionDetail, tdId, transaction.ID, product.ID, quantity, totalPrice); err != nil {
      return entities.Transaction{}, err
    }

    merchantBalance := 0
    if err := tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", merchantId).Scan(&merchantBalance); err != nil {
      return entities.Transaction{}, err
    }

    cmdTag, err := tx.Exec(ctx, "UPDATE users SET balance = $2 WHERE id = $1 AND role = 'MERCHANT';", merchantId, merchantBalance + totalPrice)
    if err != nil {
      return entities.Transaction{}, err
    }
    if cmdTag.RowsAffected() != 1 {
      return entities.Transaction{}, fmt.Errorf("No row found to update")
    }
  }

  cmdTag, err := tx.Exec(ctx, "UPDATE users SET balance = $2 WHERE id = $1 AND role = 'CUSTOMER';", userId, currentBalance)
  if err != nil {
    return entities.Transaction{}, err
  }
  if cmdTag.RowsAffected() != 1 {
    return entities.Transaction{}, fmt.Errorf("No row found to update")
  }

  if err := tx.Commit(ctx); err != nil {
    log.Println("transactions.repository: Insert.tx.Commit Err :", err)
    return entities.Transaction{}, err
  }

  return transaction, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}