package logx

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type LogMessage struct {
	Level   LogLevel
	Time    time.Time
	Message string
}

var (
	logChan chan LogMessage
	once    sync.Once
)

// ✅ Auto-start the logger when package is imported
func init() {
	once.Do(func() {
		logChan = make(chan LogMessage, 100)
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {

			log.Fatalf("Cannot open log file %v", err)
			go func() {
				for msg := range logChan {
					log.SetOutput(logFile)
					entry := fmt.Sprintf("[%s] %s: %s", msg.Time.Format(time.RFC3339), msg.Level, msg.Message)
					log.Println(entry)
				}
			}()
		}
	})
}

func Info(msg string) {
	logChan <- LogMessage{Level: INFO, Time: time.Now(), Message: msg}
}
func Warn(msg string) {
	logChan <- LogMessage{Level: WARNING, Time: time.Now(), Message: msg}
}
func Error(msg string) {
	logChan <- LogMessage{Level: ERROR, Time: time.Now(), Message: msg}
}
