package configs

import (
	"log"
	"os"
	"strconv"
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
	DATABASE_URL    string
	Token_Issue     string
	Token_Secret    string
	Token_Expire    uint64
	API_Host        string
	API_Port        string
	Timezone				string
}

type TokenConfig struct {
  IssuerName       string `json:"IssuerName"`
	JwtSignatureKey  []byte `json:"JwtSignatureKey"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

var ENV *Config

func LoadConfig() TokenConfig {
	setupDevelopmentEnv()

	if os.Getenv("APP_ENV") == "production" {
		setupProductionEnv()
	}
	
	return TokenConfig{
	  IssuerName: ENV.Token_Issue,
  	JwtSignatureKey: []byte(ENV.Token_Secret),
  	JwtSigningMethod: jwt.SigningMethodHS256,
  	JwtExpiresTime: time.Duration(ENV.Token_Expire) * time.Minute,
	}
}

func setupDevelopmentEnv() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Could not read .env file (this is normal in production): %v", err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		log.Fatal("configs.Unmarshal Err :", err)
	}
}

func setupProductionEnv() {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		ENV.DATABASE_URL = databaseURL
		log.Println("Using DATABASE_URL for connection")
	}

	if ENV.DATABASE_URL != "" {
		log.Println("Using connection string from DATABASE_URL")
	} else {
		if ENV.DB_Host == "" {
			ENV.DB_Host = os.Getenv("DB_HOST")
		}
		if ENV.DB_Port == "" {
			ENV.DB_Port = os.Getenv("DB_PORT")
		}
		if ENV.DB_User == "" {
			ENV.DB_User = os.Getenv("DB_USER")
		}
		if ENV.DB_Password == "" {
			ENV.DB_Password = os.Getenv("DB_PASSWORD")
		}
		if ENV.DB_Name == "" {
			ENV.DB_Name = os.Getenv("DB_NAME")
		}
		if ENV.API_Port == "" {
			ENV.API_Port = os.Getenv("API_PORT")
		}
		if ENV.Timezone == "" {
			ENV.Timezone = os.Getenv("TIMEZONE")
		}
	}
	
	if ENV.Token_Secret == "" {
		ENV.Token_Secret = os.Getenv("TOKEN_SECRET")
	}
	if ENV.Token_Issue == "" {
		ENV.Token_Issue = os.Getenv("TOKEN_ISSUE")
	}
	if ENV.Token_Expire == 0 {
		tokenExpire, _ := strconv.ParseUint(os.Getenv("TOKEN_EXPIRE"), 8, 64)
		ENV.Token_Expire = tokenExpire
	}
}