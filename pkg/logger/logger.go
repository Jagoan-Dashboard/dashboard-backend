package logger

import (
	"log"
	"os"
)

// Logger levels
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

// Logger interface
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// SimpleLogger implements Logger interface using standard log package
type SimpleLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

// NewSimpleLogger creates a new simple logger
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
	}
}

func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debugLogger.Printf(msg, args...)
	} else {
		l.debugLogger.Println(msg)
	}
}

func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Printf(msg, args...)
	} else {
		l.infoLogger.Println(msg)
	}
}

func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.warnLogger.Printf(msg, args...)
	} else {
		l.warnLogger.Println(msg)
	}
}

func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.errorLogger.Printf(msg, args...)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Global logger instance
var GlobalLogger Logger = NewSimpleLogger()

// Package-level functions for convenience
func Debug(msg string, args ...interface{}) {
	GlobalLogger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	GlobalLogger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	GlobalLogger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	GlobalLogger.Error(msg, args...)
}

// SetLogger sets the global logger instance
func SetLogger(logger Logger) {
	GlobalLogger = logger
}