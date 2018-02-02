package logger

import (
	log "github.com/sirupsen/logrus"
)

type NullWriter struct {
}

func (w *NullWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type JsonLogger struct {
	logger *log.Logger
}

func NewJsonLogger(hidden bool) *JsonLogger {
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

func (logger *JsonLogger) WithFields(fields map[string]interface{}) *log.Entry {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}
	return logger.logger.WithFields(fields)
}
