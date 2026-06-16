package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger creates a new logger for the application
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
