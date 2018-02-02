package logger

import (
	log "github.com/sirupsen/logrus"
)

// NullWriter does not write log
type NullWriter struct {
}

// Write write nothging
func (w *NullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

// JSONLogger is a logger JSON formatter wrap
type JSONLogger struct {
	logger *log.Logger
}

// NewJSONLogger returns a JsonLogger
func NewJSONLogger(hidden bool) *JSONLogger {
	logger := &Logger{}
	logger.logger = log.New()

	if hidden == true {
		logger.logger.Out = &NullWriter{}
	}

	logger.logger.Formatter = &log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "@timestamp",
			log.FieldKeyLevel: "@level",
			log.FieldKeyMsg:   "@message",
		},
	}

	return logger
}

// WithFields returns a log entry
func (logger *JSONLogger) WithFields(fields map[string]interface{}) *log.Entry {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}
	return logger.logger.WithFields(fields)
}
