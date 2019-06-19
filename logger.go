package apollo

import (
	logger "gopkg.in/logger.v1"
)

//Logger interface
type Logger interface {
	Warnf(format string, v ...interface{})
	Warn(v ...interface{})
	Errorf(format string, v ...interface{})
	Error(v ...interface{})
	Infof(format string, v ...interface{})
	Info(v ...interface{})
	Debugf(format string, v ...interface{})
	Debug(v ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

var log Logger

//SetLogger set user custome logger
func SetLogger(userLog Logger) {
	log = userLog
}

func setDefaultLogger() {
	log = logger.Std
}
