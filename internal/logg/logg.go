// logg uses standard log package with customized formatting and different log levels
//
// Debug and Info logger are disabled by default. To enable `EnableDebugLogger()` or `EnableInfoLogger`.
//
// Error logger is enabled by default. To disable use `DisableErrorLogger()`.
package logg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
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

// Infof is for logs like "server started".
func Infof(format string, v ...any) {
	if !infoLoggerEnabled {
		return
	}
	infoLogger.Output(2, fmt.Sprintf(format, v...))
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

// Errfo logs with internal detailed information with custom calldepth for log.Output.
// Used for logging information that happened outside this scope.
//
// calldepth = 2: This function.
//
// calldepth = 3: Outer function that called this.
//
// calldepth = ...
func Errfo(calldepth int, format string, v ...any) {
	if !errorLoggerEnabled {
		return
	}
	errorLogger.Output(calldepth, fmt.Sprintf(format, v...))
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

// Debugfo logs with internal detailed information with custom calldepth for log.Output.
// Used for logging information that happened outside this scope.
//
// calldepth = 2: This function.
//
// calldepth = 3: Outer function that called this.
//
// calldepth = ...
func Debugfo(calldepth int, format string, v ...any) {
	if !debugLoggerEnabled {
		return
	}
	debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func DebugLogger() *log.Logger {
	return debugLogger
}

// Errorf works similar to fmt.Errorf but adds line number and function name to output.
func Errorf(format string, a ...any) error {
	pc, filename, line, _ := runtime.Caller(1)

	fullFuncName := runtime.FuncForPC(pc).Name()
	// Debug(fullFuncName)
	funSplit := strings.Split(fullFuncName, "/")
	shortFuncName := funSplit[len(funSplit)-1]
	// shortFuncName = strings.Split(shortFuncName, ".")[1]

	nameSplit := strings.Split(filename, "/")
	shortFileName := nameSplit[len(nameSplit)-1]
	// Example output:
	// "file.go:69 Func: Something happened here."
	// pre := fmt.Sprintf("%s:%d\t%s", shortFileName, line, shortFuncName)
	pre := fmt.Sprintf("\n\t%s:%d [%s()] - ", shortFileName, line, shortFuncName)
	// a = append(a, "\n\t")
	// format = format + "\n\t"
	err := fmt.Errorf(pre+format, a...)
	// err = fmt.Errorf("%w\n", err)
	// aaa := fmt.Sprintf("%s", err)
	// Debug(aaa)
	return err
}

// Errorf2 works similar to fmt.Errorf2 but adds line number and function name to output.
func Errorf2(msg string, err error) error {
	pc, filename, line, _ := runtime.Caller(1)

	fullFuncName := runtime.FuncForPC(pc).Name()
	Debug(fullFuncName)
	funSplit := strings.Split(fullFuncName, "/")
	shortFuncName := funSplit[len(funSplit)-1]
	// shortFuncName = strings.Split(shortFuncName, ".")[1]

	nameSplit := strings.Split(filename, "/")
	shortFileName := nameSplit[len(nameSplit)-1]
	// Example output:
	// "file.go:69 Func: Something happened here."
	return fmt.Errorf("%s:%d\t%s: %s\n\t%w", shortFileName, line, shortFuncName, msg, err)
}

// WrapErr wraps an error with additional logg details.
func WrapErr(err error) error {
	return WrapErrWithSkip(err, 2)
}

// WrapErrWithSkip wraps error with logg details and skips stack frames.
// Used in other error wrapper functions to capture outside information.
func WrapErrWithSkip(err error, stackframes int) error {
	pc, filename, line, _ := runtime.Caller(stackframes)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funSplit := strings.Split(fullFuncName, "/")
	shortFuncName := funSplit[len(funSplit)-1]
	nameSplit := strings.Split(filename, "/")
	shortFileName := nameSplit[len(nameSplit)-1]
	err2 := fmt.Errorf("\n\t%s:%d [%s()]: %w", shortFileName, line, shortFuncName, err)
	return err2
}

// NewError creates new error with added logg details.
func NewError(text string) error {
	return WrapErrWithSkip(errors.New(text), 2)
}

// Fatal is equivalent to log.Fatal().
func Fatal(v ...any) {
	logger.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to log.Fatalf().
func Fatalf(format string, v ...any) {
	errorLogger.Output(2, "Fatal error")
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

// InfoForceOutput is for logs that must be logged no matter what config is set.
func InfoForceOutput(outputLevel int, v ...any) {
	infoLogger.Output(outputLevel, fmt.Sprint(v...))
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

// WantHave creates a well formatted string for test fail logging
//
// Calling
//
//	WantHave(http.StatusForbidden, w.Result().Status, "/login"+urlParams)
//
// Output will look like:
//
//	Want: "403"      Have: "403 Forbidden"   Info: "/login?username="
func WantHave(want any, have any, v ...any) string {
	msg := fmt.Sprintf("Want: \"%v\"\tHave: \"%v\"", want, have)

	if v == nil {
		return msg
	}

	info := "\tInfo: "
	for _, item := range v {
		info += fmt.Sprintf("\"%v\"\t", item)
	}
	msg += info
	return msg
}
