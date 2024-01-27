package configs
 
const (
  // users query
  InsertUser = "INSERT INTO users(id, name, username, email, phone_number, password, role) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, username, balance, email, phone_number, role, created_at;"
  SelectUserByEmail = "SELECT id, name, username, balance, email, phone_number, role, password FROM users WHERE email = $1;"
  SelectUserByID = "SELECT id, name, username, balance, email, phone_number, role FROM users WHERE id = $1;"
  UpdateUserBalance = "UPDATE users SET balance = $2 WHERE id = $1 RETURNING balance;"
  
  // products query
  InsertProduct = "INSERT INTO products(id, user_id, product_name, description, stock, price) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, user_id, product_name, description, stock, price;"
  SelectAllProducts = "SELECT id, user_id, product_name, description, stock, price, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2;"
  SelectAllProductsByUserID = "SELECT id, product_name, description, stock, price, created_at, updated_at FROM products WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;"
  SelectProductByID = "SELECT id, user_id, product_name, description, stock, price, created_at, updated_at FROM products WHERE id = $1 ORDER BY created_at DESC;"
  DeleteProductByID = "DELETE FROM products WHERE id = $1 AND user_id = $2 RETURNING id;"

  // transaction query
  SelectAllTransactionByUserID = "SELECT id, user_id, transaction_date, transaction_type, paid_status, cashflow, created_at, updated_at FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;"
  SelectTransactionByID = "SELECT id, transaction_date, transaction_type, paid_status, cashflow, created_at FROM transactions WHERE id = $1 AND user_id = $2;"
  InsertTransaction = "INSERT INTO transactions(id, user_id, transaction_type, paid_status, cashflow) VALUES($1, $2, $3, $4, $5) RETURNING id, user_id, transaction_date, transaction_type, paid_status, cashflow, created_at;"

  // transactions detail
  SelectAllTransactionsDetail = "SELECT id, product_id, quantity, total_price, created_at FROM transaction_details WHERE transaction_id = $1 ORDER BY created_at DESC;"
  InsertTransactionDetail = "INSERT INTO transaction_details(id, transaction_id, product_id, quantity, total_price) VALUES($1, $2, $3, $4, $5) RETURNING id, transaction_id, product_id, quantity, total_price, created_at;"
)