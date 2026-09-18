package wayfinder

import (
	"fmt"
	"os"

	"sapasora/platform/support/console"
)

type DefaultLogger struct{}

func (l *DefaultLogger) Fatal(message string, fields ...any) {
	console.Error(message, fields...)
	os.Exit(1)
}

func (l *DefaultLogger) Error(message string, fields ...any) {
	console.Error(message, fields...)
}

func (l *DefaultLogger) Warn(message string, fields ...any) {
	console.Warn(message, fields...)
}

func (l *DefaultLogger) Info(message string, fields ...any) {
	console.Info(message, fields...)
}

func (l *DefaultLogger) Debug(message string, fields ...any) {
	console.Line(fmt.Sprintf("DEBUG: "+message, fields...))
}
