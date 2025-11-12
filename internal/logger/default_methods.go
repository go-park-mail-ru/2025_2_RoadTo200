package logger

import "os"

var def_log = Logger{
	out:       os.Stdout,
	level:     TRACE,
	prefix:    "SYSTEM",
	color:     true,
	timestamp: true,
}

func Println(format string) {
	def_log.Info(format)
}
func Printf(format string, v ...interface{}) {
	def_log.Infof(format, v...)
}

func Error(format string, v ...interface{}) {
	def_log.Errorf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	def_log.Fatalf(format, v...)
}
