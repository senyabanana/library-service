package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
		DisableColors:   true,
		PadLevelText:    true,
	})

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	return &Logger{Logger: logger}
}
