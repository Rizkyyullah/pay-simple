package configs

import (
  "log"
  "time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type Config struct {
	DB_Host         string
	DB_Port         string
	DB_User         string
	DB_Password     string
	DB_Name         string
	Token_Issue     string
	Token_Secret    string
	Token_Expire    int
	API_Host        string
	API_Port        int
}

type TokenConfig struct {
  IssuerName       string `json:"IssuerName"`
	JwtSignatureKey  []byte `json:"JwtSignatureKey"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

var ENV *Config

func LoadConfig() TokenConfig {
	viper.SetConfigFile(".env")
	
	if err := viper.ReadInConfig(); err != nil {
	  log.Fatal("configs.ReadInConfig Err :", err)
	}
	
	if err := viper.Unmarshal(&ENV); err != nil {
	  log.Fatal("configs.Unmarshal Err :", err)
	}
	
	return TokenConfig{
	  IssuerName: ENV.Token_Issue,
  	JwtSignatureKey: []byte(ENV.Token_Secret),
  	JwtSigningMethod: jwt.SigningMethodHS256,
  	JwtExpiresTime: time.Duration(ENV.Token_Expire) * time.Minute,
	}
}
