package common

import (
  "context"
  "fmt"
  "log"
  
  "github.com/jackc/pgx/v5"
)

func UniqueID(conn *pgx.Conn, prefix, sequence string) (string, error) {
  currentSequence := 0
  if err := conn.QueryRow(context.Background(), "SELECT NEXTVAL($1);", sequence).Scan(&currentSequence); err != nil {
    log.Println("common.UniqueID Err :", err)
    return "", err
  }
  
  uniqueID := fmt.Sprintf("%s-%0.4d", prefix, currentSequence) // e.g. USR-0001
  return uniqueID, nil
}