package common

import (
  "fmt"
  "time"
)

func GetTimezone() *time.Location {
  location, _ := time.LoadLocation("Asia/Jakarta")
  return location
}
