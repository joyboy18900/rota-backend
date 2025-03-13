
package utils

import (
	"log"
	"os"
	"time"
)

// Logger levels
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// Logger provides a simple logging interface
type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
}

// NewLogger creates a new logger
func NewLogger() *Logger {
	return &Logger{
		debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warn:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	l.debug.Printf(format, v...)
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.info.Printf(format, v...)
}

// Warn logs a warning message
func (l *Logger) Warn(format string, v ...interface{}) {
	l.warn.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.error.Printf(format, v...)
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method, path, ip string, duration time.Duration, statusCode int) {
	l.Info("%s %s %s %v %d", method, path, ip, duration, statusCode)
}
