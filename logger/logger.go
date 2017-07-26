package logger

import (
	"bytes"
	"fmt"
	"log/syslog"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/servicekit/servicekit-go/config"
)

const (
	RequestIDKey = "request-id"
)

func insert(slice []interface{}, insertion interface{}) []interface{} {
	result := make([]interface{}, len(slice)+1)
	result[0] = insertion
	copy(result[1:], slice)
	return result
}

type LogFormatter struct {
}

func (f *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString(" ")
	b.WriteString(entry.Message)
	b.WriteString("\n")
	return b.Bytes(), nil
}

type Logger struct {
	logger *log.Logger
}

func NewLogger(serviceName, serviceVersion string, serviceENV config.ServiceENV, network, addr string, priority syslog.Priority) (*Logger, error) {
	logger := &Logger{}
	logger.logger = log.New()

	logger.logger.Formatter = &LogFormatter{}

	hook, err := logrus_syslog.NewSyslogHook(network, addr, priority, fmt.Sprintf("%s_%s_%s", serviceName, serviceVersion, serviceENV))
	if err != nil {
		return nil, err
	} else {
		logger.logger.Hooks.Add(hook)
		return logger, nil
	}
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Debugf(f, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Infof(f, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Printf(f, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Warnf(f, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Errorf(f, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Fatalf(f, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s", file, line, format)
	}

	logger.logger.Panicf(f, args...)
}

func (logger *Logger) DebugfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Debugf(f, args...)
}

func (logger *Logger) InfofWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Infof(f, args...)
}

func (logger *Logger) PrintfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Printf(f, args...)
}

func (logger *Logger) WarnfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Warnf(f, args...)
}

func (logger *Logger) ErrorfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Errorf(f, args...)
}

func (logger *Logger) FatalfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Fatalf(f, args...)
}

func (logger *Logger) PanicfWithReqID(ctx context.Context, format string, args ...interface{}) {
	f := format

	_, file, line, ok := runtime.Caller(1)
	if ok {
		f = fmt.Sprintf("%s:%d %s reqid: %v", file, line, format, ctx.Value(RequestIDKey))
	}

	logger.logger.Panicf(f, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Info(args...)
}

func (logger *Logger) Print(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Print(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		args = insert(args, fmt.Sprintf("%s:%d ", file, line))
	}

	logger.logger.Panic(args...)
}
