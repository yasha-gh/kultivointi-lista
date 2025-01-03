package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func NewLogger() *log.Logger {
	if !IsDev() {
		if cacheDir, _ := os.UserCacheDir(); cacheDir != "" {
			logFile := fmt.Sprintf("%s/kultivointi-lista/logi.log", cacheDir)
			fmt.Println("log file path", logFile)
			f, _ := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
			fmt.Println("Setting logger to file", logFile)
			// log.SetOutput(f)
			// log.SetFormatter(log.JSONFormatter) // Use JSON format
			log.SetOutput(io.MultiWriter(f, os.Stderr))
			return log.New(f)
		} else {
			fmt.Println("Logger: failed to get user cache dir", cacheDir)
		}
	}
	fmt.Println("Setting logger to stderr", "Is dev", IsDev())
	return log.New(os.Stderr)
}

func GetLogger() *log.Logger {
	if Logger == nil {
		return NewLogger()
	}
	return Logger
}
