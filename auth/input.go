package auth

type LoginInput struct {
  Email     string `json:"email"`
  Password  string `json:"password"`
}

type RegisterInput struct {
  Name              string  `json:"name" validate:"required"`
  Username          string  `json:"username" validate:"required,lowercase"`
  Email             string  `json:"email" validate:"required,email"`
  PhoneNumber       string  `json:"phoneNumber" validate:"required,number,max=13"`
  Password          string  `json:"password" validate:"required"`
  ConfirmPassword   string  `json:"confirmPassword" validate:"eqfield=Password"`
}