package common

import (
  "fmt"
  "time"

  "github.com/spf13/viper"
)

func GetTimezone() *time.Location {
  var location *time.Location
  viper.AutomaticEnv()

  envLocation := viper.GetString("TIMEZONE_IDN")
  if envLocation == "" {
    location, _ = time.LoadLocation("Asia/Jakarta")
  } else {
    location, _ = time.LoadLocation(envLocation)
  }

  fmt.Println(envLocation)

  return location
}
