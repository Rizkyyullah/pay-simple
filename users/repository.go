package users

import (
  "context"
  "log"
  "github.com/Rizkyyullah/pay-simple/configs"
  "github.com/Rizkyyullah/pay-simple/entities"
  "github.com/Rizkyyullah/pay-simple/shared/common"
  
  "github.com/jackc/pgx/v5"
)

type Repository interface {
  Insert(payload entities.User) (entities.User, error)
}

type repository struct {
  conn *pgx.Conn
}

func (r *repository) Insert(payload entities.User) (entities.User, error) {
  var user entities.User
  
  id, err := common.UniqueID(r.conn, "USR", "users_id_seq")
  if err != nil {
    log.Println("users.Repository: UniqueID Err :", err)
    return entities.User{}, err
  }
  
  if err := r.conn.QueryRow(context.Background(), configs.InsertUser,
    id,
    payload.Name,
    payload.Username,
    payload.Email,
    payload.PhoneNumber,
    payload.Password,
    payload.Role,
  ).Scan(
    &user.ID,
    &user.Name,
    &user.Username,
    &user.Balance,
    &user.Email,
    &user.PhoneNumber,
    &user.Role,
    &user.CreatedAt,
  ); err != nil {
    log.Println("users.Repository: QueryRow.Scan Err :", err)
    return entities.User{}, err
  }
  
  return user, nil
}

func NewRepository(conn *pgx.Conn) Repository {
  return &repository{conn}
}