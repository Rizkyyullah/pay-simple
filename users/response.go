package users

type GetBalanceResponse struct {
  Balance int `json:"balance,omitempty"`
}

type UserResponse struct {
  ID          string    `json:"id,omitempty"`
  Name        string    `json:"name,omitempty"`
  Username    string    `json:"username,omitempty"`
  Balance     int       `json:"balance,omitempty"`
  Email       string    `json:"email,omitempty"`
  PhoneNumber string    `json:"phoneNumber,omitempty"`
  Role        string    `json:"role,omitempty"`
  CreatedAt   string    `json:"createdAt,omitempty"`
  UpdatedAt   string    `json:"updatedAt,omitempty"`
}