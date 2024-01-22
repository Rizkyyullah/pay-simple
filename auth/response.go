package auth

type LoginResponse struct {
  Token string `json:"token"`
}

type RegisterResponse struct {
  ID          string    `json:"id"`
  Name        string    `json:"name,omitempty"`
  Username    string    `json:"username,omitempty"`
  Balance     int       `json:"balance,omitempty"`
  Email       string    `json:"email,omitempty"`
  PhoneNumber string    `json:"phoneNumber,omitempty"`
  Role        string    `json:"role,omitempty"`
}