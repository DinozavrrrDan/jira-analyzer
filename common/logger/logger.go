package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type JiraLogger struct {
	logger  *logrus.Logger
	logFile *io.Writer
	errFile *io.Writer
}

type LogLevel int

const (
	DEBUG   LogLevel = 0
	INFO    LogLevel = 1
	WARNING LogLevel = 2
	ERROR   LogLevel = 3
)

func CreateNewLogger() *JiraLogger {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel) //Trace level - самый объемный по информации

	logs, _ := os.OpenFile("././logs/logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	errors, _ := os.OpenFile("././logs/err_logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	logFile := io.MultiWriter(logs)
	errFile := io.MultiWriter(os.Stdout, errors)

	return &JiraLogger{
		logger:  logger,
		logFile: &logFile,
		errFile: &errFile,
	}
}

func (JLogger *JiraLogger) Log(logLevel LogLevel, logMessage string) {
	JLogger.logger.Out = *JLogger.logFile
	if logLevel == DEBUG {
		JLogger.logger.Debug(logMessage)
	} else if logLevel == INFO {
		JLogger.logger.Infoln(logMessage)
	} else if logLevel == WARNING {
		JLogger.logger.Warning(logMessage)
		JLogger.logger.Out = *JLogger.errFile
		JLogger.logger.Warning(logMessage)
	} else if logLevel == ERROR {
		JLogger.logger.Error(logLevel)
		JLogger.logger.Out = *JLogger.errFile
		JLogger.logger.Error(logLevel)
	}
}
