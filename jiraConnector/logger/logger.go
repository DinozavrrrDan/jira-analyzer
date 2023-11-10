package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type JiraLogger struct {
	logger  *logrus.Logger
	logFile *io.Writer
	errFile *io.Writer
}

func CreateNewLogger() *JiraLogger {
	logger := logrus.New()

	logger.SetLevel(logrus.TraceLevel) //Trace level - самый объемный по информации

	//Тут будут файлы, посоветуюсь с вами же добавлю
	//logs, _ := os.OpenFile()
	//errors, _ := os.OpenFile()

	//logFile := io.Writer(logs)
	//errFile := io.Writer(err) может добавить доп вывод в os.Stdout??

	return &JiraLogger{
		logger:  logger,
		logFile: &logFile,
		errFile: &errFile,
	}
}

// напишу после решения вопросов свыше
// func Log(logLvevel LogLevel, logMessage string)
