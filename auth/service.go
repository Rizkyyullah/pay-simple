package auth

import (
	"fmt"
	"log"
	"github.com/Rizkyyullah/pay-simple/configs"
	"github.com/Rizkyyullah/pay-simple/entities"
	"github.com/Rizkyyullah/pay-simple/shared/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user entities.User) (LoginResponse, error)
}

type jwtService struct {
	config configs.TokenConfig
}

func (j *jwtService) CreateToken(user entities.User) (LoginResponse, error) {
	claims := models.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: j.config.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.JwtExpiresTime)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserId: user.ID,
		Role: user.Role,
	}

	token := jwt.NewWithClaims(j.config.JwtSigningMethod, claims)
	ss, err := token.SignedString(j.config.JwtSignatureKey)
	if err != nil{
	  log.Println("auth.service.CreateToke Err :", err)
		return LoginResponse{}, fmt.Errorf("Oops, failed to create token")
	}

	return LoginResponse{ss}, nil
}

func NewJwtService(configs configs.TokenConfig) JwtService {
	return &jwtService{configs}
}