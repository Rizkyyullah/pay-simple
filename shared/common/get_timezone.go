package common

import (
  "time"

  "github.com/spf13/viper"
)

func GetTimezone() *time.Location {
  viper.AutomaticEnv()

  location, _ := time.LoadLocation(viper.GetString("TIMEZONE"))
  return location
}