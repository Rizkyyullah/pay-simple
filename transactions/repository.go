package transactions

import (
  "context"
  "errors"
  "fmt"
  "log"
  "math"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  "github.com/Rizkyyullah/pay-simple/shared/models"
  "github.com/Rizkyyullah/pay-simple/users"
  
  "github.com/jackc/pgx/v5"
)

var ctx = context.Background()

type Repository interface {
  FindAllByUserID(userId string, page, size int) ([]entities.Transaction, models.Paging, error)
  FindByID(id, userId string) (entities.Transaction, error)
  Insert(dto TransactionDTO) error
  UpdateBalance(dto TransactionDTO) (users.GetBalanceResponse, error)
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

func (r *repository) Insert(dto TransactionDTO) error {
  var transaction entities.Transaction

  tx, err := r.conn.Begin(ctx)
  if err != nil {
    return err
  }
  defer tx.Rollback(ctx)
  
  if err := r.insertTransaction(tx.Conn(), tx, &transaction, dto.UserID, dto.TransactionType, dto.PaidStatus, dto.Cashflow); err != nil {
    return err
  }

  currentBalance := dto.Balance
  for _, product := range dto.Products {
    productId := product.ID
    quantity := dto.Quantities[productId]
    totalPrice := product.Price * quantity
    merchantId := dto.Merchants[product.ID].ID
    currentBalance -= totalPrice
    if currentBalance < 0 {
      return fmt.Errorf("Your balance is insufficient")
    }

    if err := r.insertTransactionDetails(tx.Conn(), tx, transaction.ID, productId, quantity, totalPrice); err != nil {
      return err
    }

    merchantBalance := 0
    if err := r.updateMerchantBalance(tx, &merchantBalance, merchantId, totalPrice); err != nil {
      return err
    }
  }

  if err := r.updateCustomerBalance(tx, &currentBalance, dto.UserID); err != nil {
    return nil
  }

  if err := tx.Commit(ctx); err != nil {
    log.Println("transactions.repository: Insert.tx.Commit Err :", err)
    return err
  }

  return nil
}

func (r *repository) UpdateBalance(dto TransactionDTO) (users.GetBalanceResponse, error) {
  var balance users.GetBalanceResponse

  tx, err := r.conn.Begin(ctx)
  if err != nil {
    return users.GetBalanceResponse{}, err
  }
  defer tx.Rollback(ctx)
  
  if err := tx.QueryRow(ctx, configs.UpdateUserBalance, dto.UserID, dto.Balance + dto.Amount).Scan(&balance.Balance); err != nil {
    log.Println("transactions.repository: UpdateBalance.QueryRow.Scan Err:", err)
    return users.GetBalanceResponse{}, err
  }

  txId, _ := common.UniqueID(tx.Conn(), "TRX", "transactions_id_seq")
  cmdTag, err := tx.Exec(ctx, configs.InsertTransaction, txId, dto.UserID, dto.TransactionType, dto.PaidStatus, dto.Cashflow)
  if err != nil {
    log.Println("transactions.repository: UpdateBalance.Exec Err :", err)
    return users.GetBalanceResponse{}, err
  }

  if cmdTag.RowsAffected() != 1 {
    return users.GetBalanceResponse{}, fmt.Errorf("Failed to insert transaction")
  }

  if err := tx.Commit(ctx); err != nil {
    log.Println("transactions.repository: UpdateBalance.tx.Commit Err :", err)
    return users.GetBalanceResponse{}, err
  }

  return balance, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}

func (r *repository) insertTransaction(conn *pgx.Conn, tx pgx.Tx, transaction *entities.Transaction, userId, transactionType string, paidStatus bool, cashflow string) error {
  trxId, _ := common.UniqueID(conn, "TRX", "transactions_id_seq")
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
    return err
  }

  return nil
}

func (r *repository) insertTransactionDetails(conn *pgx.Conn, tx pgx.Tx, transactionId, productId string, quantity, totalPrice int) error {
  tdId, _ := common.UniqueID(conn, "TRXD", "transaction_details_id_seq")
  if _, err := tx.Exec(ctx, configs.InsertTransactionDetail, tdId, transactionId, productId, quantity, totalPrice); err != nil {
    return err
  }

  return nil
}

func (r *repository) updateMerchantBalance(tx pgx.Tx, merchantBalance *int, merchantId string, totalPrice int) error {
  if err := tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1", merchantId).Scan(&merchantBalance); err != nil {
    return err
  }

  cmdTag, err := tx.Exec(ctx, "UPDATE users SET balance = $2 WHERE id = $1 AND role = 'MERCHANT';", merchantId, *merchantBalance + totalPrice)
  if err != nil {
    return err
  }

  if cmdTag.RowsAffected() != 1 {
    return fmt.Errorf("No row found to update")
  }

  return nil
}

func (r *repository) updateCustomerBalance(tx pgx.Tx, customerBalance *int, customerId string) error {
  cmdTag, err := tx.Exec(ctx, "UPDATE users SET balance = $2 WHERE id = $1 AND role = 'CUSTOMER';", customerId, *customerBalance)
  if err != nil {
    return err
  }

  if cmdTag.RowsAffected() != 1 {
    return fmt.Errorf("No row found to update")
  }

  return nil
}
