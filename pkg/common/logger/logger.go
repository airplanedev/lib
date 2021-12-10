package logger

import "github.com/fatih/color"

type Logger interface {
	Log(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Step(msg string, args ...interface{})
	Suggest(title, command string, args ...interface{})
	Debug(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type NoopLogger struct{}

var _ Logger = &NoopLogger{}

func (l *NoopLogger) Log(string, ...interface{}) {
}

func (l *NoopLogger) Debug(string, ...interface{}) {
}

func (l *NoopLogger) Warning(string, ...interface{}) {
}

func (l *NoopLogger) Step(string, ...interface{}) {
}
func (l *NoopLogger) Error(string, ...interface{}) {
}
func (l *NoopLogger) Suggest(title, command string, args ...interface{}) {
}

var (
	Gray   = color.New(color.FgHiBlack).SprintfFunc()
	Blue   = color.New(color.FgHiBlue).SprintfFunc()
	Red    = color.New(color.FgHiRed).SprintfFunc()
	Yellow = color.New(color.FgHiYellow).SprintfFunc()
	Green  = color.New(color.FgHiGreen).SprintfFunc()
	Bold   = color.New(color.Bold).SprintfFunc()
)
