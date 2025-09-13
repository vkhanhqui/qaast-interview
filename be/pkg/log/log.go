package log

import (
	"be/pkg/errors"
	"fmt"
	"log"
	"os"
	"strings"

	"go.uber.org/zap"
)

type Logger interface {
	// Error will also send to Sentry if available
	Error(msg string, err error, args ...interface{})
	// Info logs at info level
	Info(msg string, args ...interface{})
	// Debug logs at debug level
	Debug(msg string, args ...interface{})

	// With adds structured context to the logger
	With(args ...interface{}) Logger
}

// PanicRecover captures the panic value, log it as error and then exit with code 1
func PanicRecover(logger Logger) {
	r := recover()
	if r == nil {
		return
	}

	err, ok := r.(error)
	if err != nil && ok {
		logger.Error("Panic with error", errors.WithStack(err))
		os.Exit(1)
		return
	}

	logger.Error("Panic occurred!", errors.Errorf("%s", r))
	os.Exit(1)
}

type LogLevel int

const (
	ErrorLevel LogLevel = 1
	InfoLevel  LogLevel = 2
	DebugLevel LogLevel = 3
)

func NewZapLogger() *ZapLogger {
	cf := zap.NewProductionConfig()
	cf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cf.DisableCaller = true
	cf.DisableStacktrace = true
	cf.EncoderConfig.TimeKey = ""
	cf.EncoderConfig.MessageKey = "message"

	logger, _ := cf.Build()

	return &ZapLogger{
		sugarZap: logger.Sugar(),
	}
}

// global logger
var gl = NewZapLogger()

// GlobalLogger returns the global logger instance
func GlobalLogger() Logger {
	return gl
}

// Error logs at error level using global logger
func Error(msg string, err error, args ...interface{}) {
	gl.Error(msg, err, args...)
}

// Info logs at info level using global logger
func Info(msg string, args ...interface{}) {
	gl.Info(msg, args...)
}

// Print outputs a simple log line
func Print(args ...interface{}) {
	log.Print(args...)
}

func Printf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

type ZapLogger struct {
	sugarZap *zap.SugaredLogger
	debug    bool
}

func (zl *ZapLogger) Error(msg string, err error, args ...interface{}) {
	zl.log(ErrorLevel, msg, err, args...)
}

func (zl *ZapLogger) Info(msg string, args ...interface{}) {
	zl.log(InfoLevel, msg, nil, args...)
}

func (zl *ZapLogger) Debug(msg string, args ...interface{}) {
	if !zl.debug {
		return
	}
	zl.log(DebugLevel, msg, nil, args...)
}

func (zl *ZapLogger) With(args ...interface{}) Logger {
	return &ZapLogger{
		sugarZap: zl.sugarZap.With(args...),
	}
}

func (zl *ZapLogger) EnableDebugLog() {
	zl.debug = true
}

func (zl *ZapLogger) log(level LogLevel, msg string, err error, args ...interface{}) {
	kvp := []interface{}{}

	ln := len(args)
	switch {
	case ln == 1:
		kvp = append(kvp, "data", args[0])
	case ln > 1:
		kvp = append(kvp, "data", args)
	}

	switch level {
	case ErrorLevel:
		str := fmt.Sprintf("%+v", err)
		str = strings.ReplaceAll(str, "\t", "    ")
		strs := strings.Split(str, "\n")

		kvp = append(kvp, "stacktrace", strs) // TODO check and add stacktrace (WithStack) if missing
		zl.sugarZap.Errorw(msg, kvp...)

	case InfoLevel:
		zl.sugarZap.Infow(msg, kvp...)

	case DebugLevel:
		zl.sugarZap.Debugw(msg, kvp...)
	}

}
