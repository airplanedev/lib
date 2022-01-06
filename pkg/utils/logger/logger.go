package logger

type Logger interface {
	Log(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}
