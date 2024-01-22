package configs

import (
  "log"

	"github.com/spf13/viper"
)

type Config struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Driver    string
	APIHost   string
	APIPort   int
}

var ENV *Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	
	if err := viper.ReadInConfig(); err != nil {
	  log.Fatal("configs.ReadInConfig Err :", err)
	}
	
	if err := viper.Unmarshal(&ENV); err != nil {
	  log.Fatal("configs.Unmarshal Err :", err)
	}
}
