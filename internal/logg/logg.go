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

const (
	Reset  = "\033[0m"
	Red    = "\033[38;5;9m"
	Green  = "\033[38;5;10m"
	Yellow = "\033[38;5;226m"
	blue   = "\033[38;5;117m"
	Gray   = "\033[38;5;245m"
	bggrey = "\033[48;5;243m"
	bgred  = "\033[48;5;217m"
)

var logger = log.New(os.Stdout, "", log.Ltime|log.Lshortfile)
var errorLogger = log.New(os.Stderr, Red+"ERROR:\t"+Gray, log.Ltime|log.Lshortfile)
var fatalLogger = log.New(os.Stderr, Red+"FATAL:\t"+Gray, log.Ltime|log.Lshortfile)
var debugLogger = log.New(os.Stdout, Green+"DEBUG:\t"+Gray, log.Ltime|log.Lshortfile)
var infoLogger = log.New(os.Stderr, blue+"INFO:\t"+Gray, log.Ltime|log.Lshortfile)
var warningLogger = log.New(os.Stdout, Yellow+"WARNING:\t"+Yellow, log.Ltime|log.Lshortfile)

var debugLoggerEnabled = false
var infoLoggerEnabled = false
var errorLoggerEnabled = true

// DefaultLogger grants access to log package.
func DefaultLogger() *log.Logger {
	return logger
}

func DebugLogger() *log.Logger {
	return debugLogger
}

func ErrorLogger() *log.Logger {
	return errorLogger
}

func InfoLogger() *log.Logger {
	return infoLogger
}

// Info is for logs like "server started".
func Info(v ...any) {
	if !infoLoggerEnabled {
		return
	}
	Alog(infoLogger, 3, "%s", v...)
}

// Warning is for logs like "there is an Error that need to be fixed but not a crucial one".
func Warning(v ...any) {
	Alog(warningLogger, 3, "%s", v...)
}

// Warningf is for logs like "there is an Error that needs to be fixed but is not a crucial one".
func Warningf(format string, v ...any) {
	// Check if the last argument is an error
	var err error
	if len(v) > 0 {
		if e, ok := v[len(v)-1].(error); ok {
			err = e
			v = v[:len(v)-1] // Remove the error from arguments
		}
	}

	// Log the formatted message
	Alog(warningLogger, 3, format, v...)

	// Log the error at the end if present
	if err != nil {
		Alog(warningLogger, 3, "Error: %s", err.Error())
	}
}

// Infof is for logs like "server started".
func Infof(format string, v ...any) {
	if !infoLoggerEnabled {
		return
	}
	Alog(infoLogger, 3, format, v...)
}

func Alog(logger *log.Logger, level int, format string, v ...any) {
	logger.Output(level, fmt.Sprintf(Reset+format, v...))
}

// Err is for logging errors which indicate internal problems (not user errors).
func Err(v ...any) {
	if !errorLoggerEnabled {
		return
	}
	Alog(errorLogger, 3, fmt.Sprint(v...))
}

// Errf is for logging errors which indicate internal problems (not user errors) with formatting.
func Errf(format string, v ...any) {
	if !errorLoggerEnabled {
		return
	}
	Alog(errorLogger, 3, format, v...)
}

// Debug is for logs with internal detailed information.
func Debug(v ...any) {
	if !debugLoggerEnabled {
		return
	}
	Alog(debugLogger, 3, fmt.Sprint(v...))
}

// Debug is for logs with internal detailed information with formatting options.
func Debugf(format string, v ...any) {
	if !debugLoggerEnabled {
		return
	}
	Alog(debugLogger, 3, format, v...)
}

// Errorf works similar to fmt.Errorf but adds line number and function name to output.
func Errorf(format string, a ...any) error {
	pre := logPrefix(2)
	err := fmt.Errorf(pre+Red+format+Reset, a...)
	return err
}

// WrapErr wraps an error with additional logg details.
func WrapErr(err error) error {
	return WrapErrWithSkip(err, 2)
}

func logPrefix(stackframes int) string {
	pc, filename, line, _ := runtime.Caller(stackframes)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funSplit := strings.Split(fullFuncName, "/")
	shortFuncName := funSplit[len(funSplit)-1]
	nameSplit := strings.Split(filename, "/")
	shortFileName := nameSplit[len(nameSplit)-1]
	return fmt.Sprintf("\n\t%s%s:%d%s [%s%s()%s] ", blue, shortFileName, line, Reset, Yellow, shortFuncName, Reset)
}

// WrapErrWithSkip wraps error with logg details and skips stack frames.
// Used in other error wrapper functions to capture outside information.
func WrapErrWithSkip(err error, stackframes int) error {
	pre := logPrefix(stackframes + 1)
	errReturn := fmt.Errorf("%s%s%w%s", pre, Red, err, Reset)
	return errReturn
}

// NewError creates new error with added logg details.
func NewError(text string) error {
	return WrapErrWithSkip(errors.New(text), 2)
}

// Fatal is equivalent to log.Fatal().
func Fatal(v ...any) {
	pre := logPrefix(3)
	Alog(fatalLogger, 3, Red+"Fatal error"+pre+Red+"%s"+Reset, v...)
	os.Exit(1)
}

// Fatalf is equivalent to log.Fatalf().
func Fatalf(format string, v ...any) {
	pre := logPrefix(3)
	Alog(fatalLogger, 3, Red+"Fatal error"+pre+Red+format+Reset, v...)
	os.Exit(1)
}

// Debug logging

// Debug logger is disabled by defaul.
func EnableDebugLogger() {
	debugLoggerEnabled = true
	Alog(debugLogger, 3, fmt.Sprint("Enabled Debug Logger"))
}

// Silent version of EnableDebugLogger Will not show that it is enabled in the logs.
func EnableDebugLoggerS() {
	debugLoggerEnabled = true
}

// Debug logger is disabled by default.
func DisableDebugLogger() {
	if debugLoggerEnabled {
		Alog(debugLogger, 3, fmt.Sprint("Disabled Debug Logger"))
	}
	debugLoggerEnabled = false
}

// Silent version of DisableDebugLogger Will not show that it is disabled in the logs.
func DisableDebugLoggerS() {
	debugLoggerEnabled = false
}

// Debug logger is disabled by default.
func DebugLoggerEnabled() bool {
	return debugLoggerEnabled
}

// Info logging

// Info logger is disabled by default.
func EnableInfoLogger() {
	infoLoggerEnabled = true
	if debugLoggerEnabled {
		Alog(debugLogger, 3, fmt.Sprint("Enabled Info Logger"))
	}
}

// Silent version of EnableInfoLogger. Will not show that it is enabled in the logs.
func EnableInfoLoggerS() {
	infoLoggerEnabled = true
}

// Info logger is disabled by default.
func DisableInfoLogger() {
	infoLoggerEnabled = false
	if debugLoggerEnabled {
		Alog(debugLogger, 3, fmt.Sprint("Disabled Info Logger"))
	}
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
	Alog(infoLogger, outputLevel, fmt.Sprint(v...))
}

// Error logging

// Error logger is enabled by default.
func EnableErrorLogger() {
	errorLoggerEnabled = true
	if debugLoggerEnabled {
		Alog(debugLogger, 3, fmt.Sprint("Enabled Error Logger"))
	}
}

// Silent version of EnableErrorLogger. Will not show that it is enabled in the logs.
func EnableErrorLoggerS() {
	errorLoggerEnabled = true
}

// Error logger is enabled by default.
func DisableErrorLogger() {
	errorLoggerEnabled = false
	if debugLoggerEnabled {
		Alog(debugLogger, 3, fmt.Sprint("Disabled Error Logger"))
	}
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
