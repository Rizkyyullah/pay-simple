package middlewares

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
	"github.com/Rizkyyullah/pay-simple/shared/common"

	"github.com/gin-gonic/gin"
)

func init() {
	readLogsFromFile()
}

type LogEntry struct {
	Timestamp   string  `json:"timestamp"`
	UserID      string  `json:"userID,omitempty"`
	Email       string  `json:"email,omitempty"`
	Action      string  `json:"action"`
	StatusCode  int     `json:"status"`
	Endpoint    string  `json:"endpoint"`
}

var mu sync.Mutex
var logEntries []LogEntry

type LogMiddleware interface {
  ActivityLogs() gin.HandlerFunc
}

type logMiddleware struct{}

func (m *logMiddleware) ActivityLogs() gin.HandlerFunc {
  return func(ctx *gin.Context) {
		ctx.Next()
		
		timestamp := time.Now().In(common.GetTimezone()).Format(time.RFC850)
		endpoint := ctx.FullPath()

		logEntry := LogEntry{
			Timestamp: timestamp,
			Action: ctx.Request.Method,
			StatusCode: ctx.Writer.Status(),
			Endpoint:  endpoint,
		}

		if userId, exists := ctx.Get("userId"); exists {
      logEntry.UserID = userId.(string)
    } else if email, exists := ctx.Get("email"); exists {
      logEntry.Email = email.(string)
    } else {
      logEntry.UserID = "Anonymous"
      logEntry.Email = "Anonymous"
    }

		// Lakukan operasi IO (menyimpan log ke dalam file) secara aman dengan menggunakan Mutex
		mu.Lock()
		defer mu.Unlock()

		logEntries = append(logEntries, logEntry)
		m.writeLogsToFile()
  }
}

func (m *logMiddleware) writeLogsToFile() {
	jsonData, err := json.MarshalIndent(logEntries, "", "  ")
	if err != nil {
		fmt.Println("Error saat marshal data log:", err)
		return
	}

	if err = ioutil.WriteFile("history.json", jsonData, 0644); err != nil {
		fmt.Println("Error saat menyimpan log ke dalam file:", err)
		return
	}

	fmt.Println("Log history telah disimpan ke dalam file: history.json")
}

func NewLogMiddleware() LogMiddleware {
  return &logMiddleware{}
}


func readLogsFromFile() {
	fileData, err := ioutil.ReadFile("history.json")
	if err != nil {
		fmt.Println("Error saat membaca file history:", err)
		return
	}

	err = json.Unmarshal(fileData, &logEntries)
	if err != nil {
		fmt.Println("Error saat unmarshal data log:", err)
		return
	}

	fmt.Println("Log history berhasil dibaca dari file history.json")
}