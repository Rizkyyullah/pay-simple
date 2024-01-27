package models

import "github.com/golang-jwt/jwt/v5"

type MyCustomClaims struct {
	jwt.RegisteredClaims
	UserID      string  `json:"userId"`
	Role        string  `json:"role"`
	Authorized  bool    `json:"authorized"`
}
