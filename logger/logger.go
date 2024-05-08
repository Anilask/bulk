package logger

// ILogger the interface provide the methods to send log data by type
type ILogger interface {
	Debug(mgs string)
	Info(mgs string)
	Warn(mgs string)
	Error(mgs string)
	Fatal(mgs string)
	Panic(mgs string)

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Infow(fields map[string]interface{}, message string)
}
