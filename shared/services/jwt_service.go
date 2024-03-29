package services

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
	CreateToken(user entities.User) (string, error)
	ParseToken(tokenHeader string) (jwt.MapClaims, error)
}

type jwtService struct {
	config configs.TokenConfig
}

func (j *jwtService) CreateToken(user entities.User) (string, error) {
	claims := models.MyCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: j.config.IssuerName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.config.JwtExpiresTime)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
		UserID: user.ID,
		Role: user.Role,
		Authorized: true,
	}

	token := jwt.NewWithClaims(j.config.JwtSigningMethod, claims)
	ss, err := token.SignedString(j.config.JwtSignatureKey)
	if err != nil{
	  log.Println("auth.service.CreateToke Err :", err)
		return "", fmt.Errorf("failed to create token")
	}

	return ss, nil
}

func (j *jwtService) ParseToken(tokenHeader string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenHeader, func(token *jwt.Token) (interface{}, error) {
		return j.config.JwtSignatureKey, nil
	})

	if err != nil {
	  log.Println("auth.service.ParseToken Err :", err)
		return nil, fmt.Errorf("Failed to verify token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Failed to claim token")
	}

	return claims, nil
}

func NewJwtService(config configs.TokenConfig) JwtService {
	return &jwtService{config}
}