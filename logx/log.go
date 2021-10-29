package logx

import (
	clogx "github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go-common/util/fsx"
	"os"
)

var logDir string
var loggers = map[string]clogx.Logger{}
var runtimeLogger clogx.Logger
var requestLogLogger clogx.Logger
var executeTimeLogLogger clogx.Logger

func WithLogDir(dir string) {
	dir = fsx.GetRealpath(dir)
	
	if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
		logDir = dir
	}
}

func GetLogDir() string {
	return logDir
}

func WithLogger(name string, logger clogx.Logger) {
	loggers[name] = logger
}

func Channel(name string) clogx.Logger {
	logger := loggers[name]
	
	if logger == nil {
		logger = NewNoopLogger()
	}
	
	return logger
}

func WithRuntimeLogger(logger clogx.Logger) {
	runtimeLogger = logger
}

func GetRuntimeLogger() clogx.Logger {
	logger := runtimeLogger
	
	if logger == nil {
		logger = NewNoopLogger()
	}
	
	return logger
}

func WithRequestLogLogger(logger clogx.Logger) {
	requestLogLogger = logger
}

func RequestLogEnabled() bool {
	return requestLogLogger != nil
}

func GetRequestLogLogger() clogx.Logger {
	logger := requestLogLogger

	if logger == nil {
		logger = NewNoopLogger()
	}

	return logger
}

func WithExecuteTimeLogLogger(logger clogx.Logger) {
	executeTimeLogLogger = logger
}

func ExecuteTimeLogEnabled() bool {
	return executeTimeLogLogger != nil
}

func GetExecuteTimeLogLogger() clogx.Logger {
	logger := executeTimeLogLogger

	if logger == nil {
		logger = NewNoopLogger()
	}

	return logger
}

func Log(level interface{}, args ...interface{}) {
	GetRuntimeLogger().Log(level, args...)
}

func Logf(level interface{}, format string, args ...interface{}) {
	GetRuntimeLogger().Logf(level, format, args...)
}

func Trace(args ...interface{}) {
	GetRuntimeLogger().Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	GetRuntimeLogger().Tracef(format, args...)
}

func Debug(args ...interface{}) {
	GetRuntimeLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetRuntimeLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	GetRuntimeLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	GetRuntimeLogger().Infof(format, args...)
}

func Warn(args ...interface{}) {
	GetRuntimeLogger().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	GetRuntimeLogger().Warnf(format, args...)
}

func Error(args ...interface{}) {
	GetRuntimeLogger().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	GetRuntimeLogger().Errorf(format, args...)
}

func Panic(args ...interface{}) {
	GetRuntimeLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	GetRuntimeLogger().Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	GetRuntimeLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	GetRuntimeLogger().Infof(format, args...)
}
