package logger

import "github.com/sirupsen/logrus"

type LoggerProvider interface {
	Error(string)
	Errorf(string, string)
}

func Error(message ...interface{}) {
	logrus.Error(message...)
}

func Errorf(format string, message ...interface{}) {
	logrus.Errorf(format, message...)
}
