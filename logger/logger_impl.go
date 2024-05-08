package logger

import (
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"bulk/constant"
)

type Logger struct {
	service string
	env     string
	log     *log.Logger
}

var Log *Logger

func New(serviceName, env, logLevel string) *Logger {
	Log = Init(serviceName, env, logLevel)
	return Log
}

func (l Logger) Debug(mgs string) {
	pc, file, line, _ := runtime.Caller(1)

	fields := l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	)
	if strings.EqualFold(l.env, constant.EnvDev) || strings.EqualFold(l.env, constant.EnvStage) {
		fields.Debug(mgs)
	} else {
		fields.Info(mgs)
	}
}

func (l Logger) Info(mgs string) {
	l.log.WithFields(
		log.Fields{
			"service": l.service,
			"env":     l.env,
		},
	).Info(mgs)
}

func (l Logger) Warn(mgs string) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Warn(mgs)
}

func (l Logger) Error(mgs string) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Error(mgs)
}

func (l Logger) Fatal(mgs string) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Fatal(mgs)
}

func (l Logger) Panic(mgs string) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Panic(mgs)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fields := l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	)
	if strings.EqualFold(l.env, constant.EnvDev) || strings.EqualFold(l.env, constant.EnvStage) {
		fields.Debugf(format, args...)
	} else {
		fields.Infof(format, args...)
	}
}

func (l Logger) Infof(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Infof(format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Warnf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Errorf(format, args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Fatalf(format, args...)
}

func (l Logger) Panicf(format string, args ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).Panicf(format, args...)
}

func (l Logger) Infow(fields map[string]interface{}, message string) {
	pc, file, line, _ := runtime.Caller(1)
	l.log.WithFields(
		log.Fields{
			"service":  l.service,
			"env":      l.env,
			"file":     file,
			"function": runtime.FuncForPC(pc).Name(),
			"line":     line,
		},
	).WithFields(fields).Info(message)
}

func Init(service string, env string, logLevel string) *Logger {
	ll := log.New()
	// Log as JSON instead of the default ASCII formatter.
	ll.SetFormatter(
		&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyTime:  "timestamp",
				log.FieldKeyLevel: "severity",
				log.FieldKeyMsg:   "message",
			},
		},
	)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	ll.SetOutput(os.Stdout)

	// Include caller method to the log trace
	ll.SetReportCaller(false)

	if logLevel == "" {
		logLevel = "error"
	}
	lv, err := log.ParseLevel(logLevel)
	if err != nil {
		lv = log.ErrorLevel
	}
	// Only log the warning severity or above.
	ll.SetLevel(lv)
	// Create a new instance of the logger. You can have any number of instances.

	return &Logger{
		service: service,
		env:     env,
		log:     ll,
	}
}
