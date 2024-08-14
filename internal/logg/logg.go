// logg uses standard log package with customized formatting and different log levels
package logg

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
var errorLogger = log.New(os.Stderr, "ERROR:\t", log.Ltime|log.Lshortfile)
var debugLogger = log.New(os.Stdout, "DEBUG:\t", log.Ltime|log.Lshortfile)
var infoLogger = log.New(os.Stderr, "INFO:\t", log.Ltime|log.Lshortfile)

// Info is for logs like "server started"
func Info(v ...any) {
	infoLogger.Output(2, fmt.Sprint(v...))
}

// Err is for logging errors which indicate internal problems (not user errors)
func Err(v ...any) {
	errorLogger.Output(2, fmt.Sprint(v...))
}

// Errf is for logging errors which indicate internal problems (not user errors) with formatting.
func Errf(format string, v ...any) {
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}

// Debug is for logs with internal detailed information.
func Debug(v ...any) {
	debugLogger.Output(2, fmt.Sprint(v...))
}

// Debug is for logs with internal detailed information with formatting options.
func Debugf(format string, v ...any) {
	debugLogger.Output(2, fmt.Sprintf(format, v...))
}

// Fatal is equivalent to log.Fatal().
func Fatal(v ...any) {
	logger.Fatal(v...)
}

// Fatalf is equivalent to log.Fatalf().
func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

// DefaultLogger grants access to log package.
func DefaultLogger() *log.Logger {
	return logger
}
