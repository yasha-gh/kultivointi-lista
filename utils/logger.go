package utils

import (
	"os"

	"github.com/charmbracelet/log"
)
var Logger *log.Logger
func NewLogger() *log.Logger {
	return log.New(os.Stderr)
}

func GetLogger() *log.Logger {
	if Logger == nil {
		return NewLogger() 
	}
	return Logger
}
