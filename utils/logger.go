package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	logger *logrus.Logger
)

type Logger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

func (logger *Logger) AddFields(fields logrus.Fields) *Logger {
	logger.fields = fields
	return logger
}

func (logger *Logger) Info(args ...interface{}) {
	logger.logger.WithFields(logger.fields).Info(args)
}

func init() {
	logger = logrus.StandardLogger()

	logger.Out = os.Stdout

	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}
}

func GetLogger() *logrus.Logger {
	return logger
}

func GetLoggerEntry() *logrus.Entry {
	return logger.WithFields(logrus.Fields{"created": time.Now().UnixNano() / 1e6})
}
