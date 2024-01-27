package auth

type LoginInput struct {
  Email     string `json:"email"`
  Password  string `json:"password"`
}

type RegisterInput struct {
  Name              string  `json:"name" validate:"required,min=3"`
  Username          string  `json:"username" validate:"required,lowercase,min=3"`
  Email             string  `json:"email" validate:"required,email"`
  PhoneNumber       string  `json:"phoneNumber" validate:"required,number,min=12,max=13"`
  Password          string  `json:"password" validate:"required,min=8"`
  ConfirmPassword   string  `json:"confirmPassword" validate:"eqfield=Password"`
}