package configs

const (
  InsertUser = "INSERT INTO users(id, name, username, email, phone_number, password, role) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, name, username, balance, email, phone_number, role, created_at;"
)