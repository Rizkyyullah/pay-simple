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
	Port        int
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
	viper.AutomaticEnv()
	
	if err := viper.ReadInConfig(); err != nil {
	   log.Println("configs.ReadInConfig Err :", err)
	  if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
  		ENV = &Config{
  		  DB_Host: viper.GetString("DB_HOST"),
      	DB_Port: viper.GetString("DB_PORT"),
      	DB_User: viper.GetString("DB_USER"),
      	DB_Password: viper.GetString("DB_PASSWORD"),
      	DB_Name: viper.GetString("DB_NAME"),
      	Token_Issue: viper.GetString("TOKEN_ISSUE"),
      	Token_Secret: viper.GetString("TOKEN_SECRET"),
      	Token_Expire: viper.GetInt("TOKEN_EXPIRE"),
      	API_Host: viper.GetString("API_HOST"),
      	Port: viper.GetInt("PORT"),
  		}
  	}
	}
	
	if err := viper.Unmarshal(&ENV); err != nil {
	  log.Println("configs.Unmarshal Err :", err)
	}
	
	return TokenConfig{
	  IssuerName: ENV.Token_Issue,
  	JwtSignatureKey: []byte(ENV.Token_Secret),
  	JwtSigningMethod: jwt.SigningMethodHS256,
  	JwtExpiresTime: time.Duration(ENV.Token_Expire) * time.Minute,
	}
}
