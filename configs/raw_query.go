package configs

const (
  // users query
  InsertUser = "INSERT INTO users(id, name, username, email, phone_number, password, role) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, username, balance, email, phone_number, role, created_at;"
  SelectUserByEmail = "SELECT id, name, username, balance, email, phone_number, role, password FROM users WHERE email = $1;"
  SelectUserByID = "SELECT id, name, username, balance, email, phone_number, role FROM users WHERE id = $1;"
  InsertProduct = "INSERT INTO products(id, user_id, product_name, description, stock, price) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, user_id, product_name, description, stock, price;"
)