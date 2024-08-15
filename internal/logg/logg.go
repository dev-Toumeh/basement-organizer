// logg uses standard log package with customized formatting and different log levels
//
// Debug and Info logger are disabled by default. To enable `EnableDebugLogger()` or `EnableInfoLogger`.
//
// Error logger is enabled by default. To disable use `DisableErrorLogger()`.
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

var debugLoggerEnabled = false
var infoLoggerEnabled = false
var errorLoggerEnabled = true

// Info is for logs like "server started".
func Info(v ...any) {
	if !infoLoggerEnabled {
		return
	}
	infoLogger.Output(2, fmt.Sprint(v...))
}

// Err is for logging errors which indicate internal problems (not user errors).
func Err(v ...any) {
	if !errorLoggerEnabled {
		return
	}
	errorLogger.Output(2, fmt.Sprint(v...))
}

// Errf is for logging errors which indicate internal problems (not user errors) with formatting.
func Errf(format string, v ...any) {
	if !errorLoggerEnabled {
		return
	}
	errorLogger.Output(2, fmt.Sprintf(format, v...))
}

// Debug is for logs with internal detailed information.
func Debug(v ...any) {
	if !debugLoggerEnabled {
		return
	}
	debugLogger.Output(2, fmt.Sprint(v...))
}

// Debug is for logs with internal detailed information with formatting options.
func Debugf(format string, v ...any) {
	if !debugLoggerEnabled {
		return
	}
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

//
// Debug logging

// Debug logger is disabled by defaul.
func EnableDebugLogger() {
	debugLoggerEnabled = true
	debugLogger.Output(2, fmt.Sprint("Enabled Debug Logger"))
}

// Silent version of EnableDebugLogger Will not show that it is enabled in the logs.
func EnableDebugLoggerS() {
	debugLoggerEnabled = true
}

// Debug logger is disabled by default.
func DisableDebugLogger() {
	debugLoggerEnabled = false
	debugLogger.Output(2, fmt.Sprint("Disabled Debug Logger"))
}

// Silent version of DisableDebugLogger Will not show that it is disabled in the logs.
func DisableDebugLoggerS() {
	debugLoggerEnabled = false
}

// Debug logger is disabled by default.
func DebugLoggerEnabled() bool {
	return debugLoggerEnabled
}

//
// Info logging

// Info logger is disabled by default.
func EnableInfoLogger() {
	infoLoggerEnabled = true
	debugLogger.Output(2, fmt.Sprint("Enabled Info Logger"))
}

// Silent version of EnableInfoLogger. Will not show that it is enabled in the logs.
func EnableInfoLoggerS() {
	infoLoggerEnabled = true
}

// Info logger is disabled by default.
func DisableInfoLogger() {
	infoLoggerEnabled = false
	debugLogger.Output(2, fmt.Sprint("Disabled Info Logger"))
}

// Silent version of DisableInfoLogger. Will not show that it is disabled in the logs.
func DisableInfoLoggerS() {
	infoLoggerEnabled = false
}

// Info logger is disabled by default.
func InfoLoggerEnabled() bool {
	return infoLoggerEnabled
}

//
// Error logging

// Error logger is enabled by default.
func EnableErrorLogger() {
	errorLoggerEnabled = true
	debugLogger.Output(2, fmt.Sprint("Enabled Error Logger"))
}

// Silent version of EnableErrorLogger. Will not show that it is enabled in the logs.
func EnableErrorLoggerS() {
	errorLoggerEnabled = true
}

// Error logger is enabled by default.
func DisableErrorLogger() {
	errorLoggerEnabled = false
	debugLogger.Output(2, fmt.Sprint("Disabled Error Logger"))
}

// Silent version of DisableErrorLogger. Will not show that it is disabled in the logs.
func DisableErrorLoggerS() {
	errorLoggerEnabled = false
}

// Error logger is enabled by default.
func ErrorLoggerEnabled() bool {
	return errorLoggerEnabled
}
