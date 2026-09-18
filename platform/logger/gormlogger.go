package logger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormLoggerImpl struct {
	*Logger

	SlowThreshold             time.Duration
	logMode                   logger.LogLevel
	IgnoreRecordNotFoundError bool
}

func (l *Logger) GetGormLogger(level logger.LogLevel) *GormLoggerImpl {
	return &GormLoggerImpl{
		Logger:                    l,
		logMode:                   level,
		SlowThreshold:             200 * time.Millisecond,
		IgnoreRecordNotFoundError: true,
	}
}

func (l *GormLoggerImpl) LogMode(level logger.LogLevel) logger.Interface {
	return &GormLoggerImpl{
		Logger:                    l.Logger,
		logMode:                   level,
		SlowThreshold:             l.SlowThreshold,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
	}
}

func (l *GormLoggerImpl) Info(ctx context.Context, msg string, data ...any) {
	if l.logMode >= logger.Info {
		l.Logger.Info(msg, data...)
	}
}

func (l *GormLoggerImpl) Warn(ctx context.Context, msg string, data ...any) {
	if l.logMode >= logger.Warn {
		l.Logger.Warn(msg, data...)
	}
}

func (l *GormLoggerImpl) Error(ctx context.Context, msg string, data ...any) {
	if l.logMode >= logger.Error {
		l.Logger.Error(msg, data...)
	}
}

func (l *GormLoggerImpl) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	if l.logMode <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logMode >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Logger.Trace("SQL query failed", "sql", sql, "error", err)
		} else {
			l.Logger.Trace("SQL query failed", "sql", sql, "error", err, "rows", rows, "time", elapsed)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logMode >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Logger.Warn(slowLog, "error", err, "sql", sql)
		} else {
			l.Logger.Warn(slowLog, "error", err, "sql", sql, "rows", rows, "time", elapsed)
		}
	case l.logMode == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Logger.Info("SQL query", "sql", sql)
		} else {
			l.Logger.Info("SQL query", "sql", sql, "rows", rows, "time", elapsed)
		}
	}
}
