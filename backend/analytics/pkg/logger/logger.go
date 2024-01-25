package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
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

func CreateNewLogger() *Logger {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel) //Trace level - самый объемный по информации
	logger.SetFormatter(&logrus.JSONFormatter{})

	logs, _ := os.OpenFile("backend/analytics/logs/logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	errors, _ := os.OpenFile("backend/analytics/logs/err_logs.log", os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	logFile := io.MultiWriter(logs)
	errFile := io.MultiWriter(os.Stdout, errors)

	return &Logger{
		logger:  logger,
		logFile: &logFile,
		errFile: &errFile,
	}
}

func (log *Logger) Log(logLevel LogLevel, logMessage string) {
	log.logger.Out = *log.logFile
	if logLevel == DEBUG {
		log.logger.Debug(logMessage)
	} else if logLevel == INFO {
		log.logger.Infoln(logMessage)
	} else if logLevel == WARNING {
		log.logger.Warning(logMessage)
		log.logger.Out = *log.errFile
		log.logger.Warning(logMessage)
		fmt.Print(logMessage)
	} else if logLevel == ERROR {
		log.logger.Error(logMessage)
		log.logger.Out = *log.errFile
		log.logger.Error(logMessage)
		fmt.Print(logMessage)
	}
}
