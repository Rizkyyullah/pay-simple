package configs

import (
  "fmt"
  "io"
  "log"
  "os"
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
	
	if os.Getenv("APPMODE") == "DEPLOY" {
	  file, err := os.Create(".env")
    if err != nil {
        log.Fatal("Error :", err)
    }
    defer file.Close()

    // teks := "Ini adalah contoh menulis ke file menggunakan bahasa Golang.\n"
    // data := []byte(teks)
    // _, err = file.Write(data)
    // if err != nil {
    //     fmt.Println("Error:", err)
    //     os.Exit(1)
    // }

    text := `
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
API_PORT=$API_PORT
PORT=$PORT
TOKEN_ISSUE=$TOKEN_ISSUE
TOKEN_SECRET=$TOKEN_SECRET
TOKEN_EXPIRE=$TOKEN_EXPIRE
TIMEZONE=$TIMEZONE
  `
    _, err = io.WriteString(file, text)
    if err != nil {
        log.Fatal("Error :", err)
    }

    // Menyimpan perubahan ke file
    err = file.Sync()
    // Jika terjadi error, tampilkan pesan error dan keluar dari program
    if err != nil {
        log.Fatal("Error :", err)
    }

    // Menampilkan pesan sukses
    fmt.Println("Sukses menulis ke file.")
	}
	
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
