package logger

import (
	"intmax2-withdrawal/internal/logger"
	"intmax2-withdrawal/internal/logger/logrus"
)

func New(logLevel, timeFormat string, logJSON, logLines bool) logger.Logger {
	return logrus.New(logLevel, timeFormat, logJSON, logLines)
}
