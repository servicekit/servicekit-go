package logger

import (
	"fmt"
	"log/syslog"
	"runtime"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/servicekit/servicekit-go/config"
)

func insert(slice []interface{}, insertion interface{}) []interface{} {
	result := make([]interface{}, len(slice)+1)
	result[0] = insertion
	copy(result[1:], slice)
	return result
}

// Logger is a abstraction base on log.Logger
type Logger struct {
	logger *log.Logger
	Active bool
}

// NewLogger returns a Logger
func NewLogger(serviceName, serviceVersion string, serviceENV config.ServiceENV, network, addr string, priority syslog.Priority) (*Logger, error) {
	logger := &Logger{Active: true}
	logger.logger = log.New()

	logger.logger.Formatter = &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "@level",
			log.FieldKeyMsg:   "@message",
		},
	}

	hook, err := logrus_syslog.NewSyslogHook(network, addr, priority, fmt.Sprintf("%s_%s_%s", serviceName, serviceVersion, serviceENV))
	if err != nil {
		return nil, err
	}

	logger.logger.Hooks.Add(hook)

	return logger, nil
}

// Debugf will invoke logrus.Debugf
func (logger *Logger) Debugf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Debugf(f, args...)
}

// Infof will invoke logrus.Infof
func (logger *Logger) Infof(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Infof(f, args...)
}

// Printf will invoke logrus.Printf
func (logger *Logger) Printf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Printf(f, args...)
}

// Warnf will invoke logrus.Warnf
func (logger *Logger) Warnf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Warnf(f, args...)
}

// Errorf will invoke logrus.Warnf
func (logger *Logger) Errorf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Errorf(f, args...)
}

// Fatalf will invoke logrus.Fatalf
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Fatalf(f, args...)
}

// Panicf will invoke logrus.Panicf
func (logger *Logger) Panicf(format string, args ...interface{}) {
	if logger.Active == false {
		return
	}

	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Panicf(f, args...)
}

// Debug will invoke logrus.Debug
func (logger *Logger) Debug(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Debug(args...)
}

// Info will invoke logrus.Info
func (logger *Logger) Info(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Info(args...)
}

// Print will invoke logrus.Print
func (logger *Logger) Print(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Print(args...)
}

// Warn will invoke logrus.Warn
func (logger *Logger) Warn(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Warn(args...)
}

// Error will invoke logrus.Error
func (logger *Logger) Error(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Error(args...)
}

// Fatal will invoke logrus.Fatal
func (logger *Logger) Fatal(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Fatal(args...)
}

// Panic will invoke logrus.Panic
func (logger *Logger) Panic(args ...interface{}) {
	if logger.Active == false {
		return
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Panic(args...)
}

// WithFields returns an Entry with fields
func (logger *Logger) WithFields(fields map[string]interface{}) *log.Entry {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}

	return logger.logger.WithFields(fields)
}
