package logg

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
var errorLogger = log.New(os.Stderr, "ERR:\t", log.Ltime|log.Lshortfile)
var debugLogger = log.New(os.Stdout, "DEBUG:\t", log.Ltime|log.Lshortfile)
var infoLogger = log.New(os.Stderr, "INFO:\t", log.Ltime|log.Lshortfile)

func Info(v ...any) {
	infoLogger.Output(2, fmt.Sprint(v...))
}

func Err(v ...any) {
	errorLogger.Output(2, fmt.Sprint(v...))
}

func Debug(v ...any) {
	debugLogger.Output(2, fmt.Sprint(v...))
}

func Debugf(format string, v ...any) {
	debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

func DefaultLogger() *log.Logger {
	return logger
}
