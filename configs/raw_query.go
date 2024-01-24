package configs

const (
  // common query
  SelectTotalRows = "SELECT COUNT(*) total_rows FROM products;"
  
  // users query
  InsertUser = "INSERT INTO users(id, name, username, email, phone_number, password, role) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, username, balance, email, phone_number, role, created_at;"
  SelectUserByEmail = "SELECT id, name, username, balance, email, phone_number, role, password FROM users WHERE email = $1;"
  SelectUserByID = "SELECT id, name, username, balance, email, phone_number, role FROM users WHERE id = $1;"

  // products query
  InsertProduct = "INSERT INTO products(id, user_id, product_name, description, stock, price) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, user_id, product_name, description, stock, price;"
  SelectAllProducts = "SELECT id, user_id, product_name, description, stock, price, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2;"
)