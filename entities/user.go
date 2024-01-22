package entities

import "time"

type User struct {
  ID          string    `json:"id"`
  Name        string    `json:"name,omitempty"`
  Username    string    `json:"username,omitempty"`
  Balance     int       `json:"balance,omitempty"`
  Email       string    `json:"email,omitempty"`
  PhoneNumber string    `json:"phoneNumber,omitempty"`
  Password    string    `json:"password,omitempty"`
  Role        string    `json:"role,omitempty"`
  CreatedAt   time.Time `json:"createdAt,omitempty"`
  UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}
