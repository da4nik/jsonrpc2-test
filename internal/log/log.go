package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// InitLogger initializes logger
func InitLogger() {
	logger.Level = logrus.DebugLevel
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

// Errorf prints error message
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Infof prints info message
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Debugf prints debug message
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
