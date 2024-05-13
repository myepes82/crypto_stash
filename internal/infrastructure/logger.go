package infrastructure

import (
	"fmt"
)

const (
	resetLogColor  = "\033[0m"
	redLogColor    = "\033[31m"
	greenLogColor  = "\033[32m"
	yellowLogColor = "\033[33m"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) LogWarm(message string) {
	log(message, yellowLogColor)
}

func (l *Logger) LogDebug(message string) {
	log(message, greenLogColor)
}

func (l *Logger) LogError(message string) {
	log(message, redLogColor)
}

func log(message string, color string) {
	fmt.Println(color + message + resetLogColor)
}
