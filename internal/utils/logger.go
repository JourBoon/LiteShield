package utils

import (
	"log"
	"time"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	log.Printf("[Info] %s : "+format, append([]interface{}{time.Now().Format("15:04:05")}, v...)...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	log.Printf("[Error] %s : "+format, append([]interface{}{time.Now().Format("15:04:05")}, v...)...)
}
