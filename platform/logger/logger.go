package logger

import (
	"fmt"
	"io"
	"os"
	"path"

	"sapasora/platform/config"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	core zerolog.Logger
}

func NewLogger(cfg config.Config) *Logger {
	var w []io.Writer

	if cfg.GetBool("app.log.console", true) {
		w = append(w, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if cfg.GetBool("app.log.file", true) {
		fw, err := rollingFile(cfg)
		if err != nil {
			panic(fmt.Sprintf("failed to configure Rolling Log File, err: %v", err))
		}
		w = append(w, fw)
	}

	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch cfg.GetString("app.log.level", "debug") {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "none":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	logger := zerolog.New(io.MultiWriter(w...)).
		With().
		Timestamp().
		Logger()

	return &Logger{
		core: logger,
	}
}

func (l *Logger) Scope(name string) *Logger {
	return &Logger{
		core: l.core.With().Str("scope", name).Logger(),
	}
}

func (l *Logger) GetLevel() string {
	switch l.core.GetLevel() {
	case zerolog.TraceLevel:
		return "trace"
	case zerolog.DebugLevel:
		return "debug"
	case zerolog.InfoLevel:
		return "info"
	case zerolog.WarnLevel:
		return "warn"
	case zerolog.ErrorLevel:
		return "error"
	case zerolog.FatalLevel:
		return "fatal"
	case zerolog.PanicLevel:
		return "panic"
	case zerolog.NoLevel:
		return "none"
	case zerolog.Disabled:
		return "disable"
	default:
		return "none"
	}
}

func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.log(l.core.Debug(), msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.log(l.core.Info(), msg, keysAndValues...)
}

func (l *Logger) Trace(msg string, keysAndValues ...any) {
	l.log(l.core.Trace(), msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...any) {
	l.log(l.core.Warn(), msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...any) {
	l.log(l.core.Error(), msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...any) {
	l.log(l.core.Fatal(), msg, keysAndValues...)
}

func (l *Logger) GetLogger() *zerolog.Logger {
	return &l.core
}

func (l *Logger) log(event *zerolog.Event, msg string, keysAndValues ...any) {
	// Add key-value pairs
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			key, val := keysAndValues[i], keysAndValues[i+1]
			if keyStr, ok := key.(string); ok {
				event = event.Interface(keyStr, redactField(keyStr, val))
			}
		}
	}
	event.Msg(msg)
}

func rollingFile(cfg config.Config) (io.Writer, error) {
	logDir := cfg.GetString("app.log.dir", "./storage/logs")
	logName := cfg.GetString("app.log.name", cfg.GetString("app.name", "application"))
	err := os.MkdirAll(logDir, 0744)
	if err != nil {
		return nil, err
	}

	return &lumberjack.Logger{
		Filename:   path.Join(logDir, logName+".log"),
		MaxBackups: cfg.GetInt("app.log.max_backup", 20), // files
		MaxSize:    cfg.GetInt("app.log.max_size", 50),   // megabytes
		MaxAge:     cfg.GetInt("app.log.max_age", 7),     // days
		Compress:   true,
	}, nil
}
