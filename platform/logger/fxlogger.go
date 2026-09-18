package logger

import (
	"strings"

	"go.uber.org/fx/fxevent"
)

type FxLoggerImpl struct {
	*Logger
}

func (l *Logger) GetFxLogger() *FxLoggerImpl {
	return &FxLoggerImpl{l}
}

func (l *FxLoggerImpl) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Trace(
			"OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Error(
				"OnStart hook failed",
				"error", e.Err.Error(),
				"callee", e.FunctionName,
				"caller", e.CallerName,
			)
		} else {
			l.Trace(
				"OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Trace(
			"OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Error(
				"OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"error", e.Err.Error(),
			)
		} else {
			l.Trace(
				"OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Error(
				"supplied",
				"error", e.Err.Error(),
				"type", e.TypeName,
				"module", e.ModuleName,
			)
		} else {
			l.Trace(
				"supplied",
				"type", e.TypeName,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Trace(
				"provided",
				"constructor", e.ConstructorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.Error(
				"error encountered while applying options",
				"error", e.Err.Error(),
				"module", e.ModuleName,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Trace(
				"decorated",
				"decorator", e.DecoratorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			l.Error(
				"error encountered while applying options",
				"error", e.Err.Error(),
				"module", e.ModuleName,
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Trace(
			"invoking",
			"function", e.FunctionName,
			"module", e.ModuleName,
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Error(
				"invoke failed",
				"error", e.Err.Error(),
				"module", e.ModuleName,
				"function", e.FunctionName,
			)
		} else {
			l.Trace(
				"invoked",
				"function", e.FunctionName,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Stopping:
		l.Trace(
			"received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Error(
				"stop failed",
				"error", e.Err.Error(),
			)
		}
	case *fxevent.RollingBack:
		l.Error(
			"start failed, rolling back",
			"error", e.StartErr,
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Error(
				"rollback failed",
				"error", e.Err.Error(),
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Error(
				"start failed",
				"error", e.Err.Error(),
			)
		} else {
			l.Trace("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Error(
				"custom logger initialization failed",
				"error", e.Err.Error(),
			)
		} else {
			l.Trace(
				"initialized custom fxevent.er",
				"function", e.ConstructorName,
			)
		}
	}
}
