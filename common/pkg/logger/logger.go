package logger

import (
	"log"
	"os"
	"time"
)

func getLogger() *log.Logger {
	day := time.Now().Format("2006-01-02")
	logFileName := "./logs/" + day + ".log"
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file ", err)
	}

	logger := log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}

func Info(args ...interface{}) {
	logger := getLogger()
	logger.SetPrefix("INFO ")
	logger.Println(args...)
}

// Danger 为什么不命名为 error？避免和 error 类型重名
func Danger(args ...interface{}) {
	logger := getLogger()
	logger.SetPrefix("ERROR ")
	logger.Println(args...)
}

func Warning(args ...interface{}) {
	logger := getLogger()
	logger.SetPrefix("WARNING ")
	logger.Println(args...)
}
