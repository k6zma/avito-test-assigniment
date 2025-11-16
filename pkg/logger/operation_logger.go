package logger

import (
	"log/slog"
	"time"
)

type OperationLogger struct {
	logger *slog.Logger
	start  time.Time
	msg    string
	attrs  []any
}

func StartOperation(log *slog.Logger, msg string, attrs ...any) *OperationLogger {
	log.Debug("Started "+msg, attrs...)

	return &OperationLogger{
		logger: log,
		start:  time.Now(),
		msg:    msg,
		attrs:  attrs,
	}
}

func (r *OperationLogger) FinishOperation(err *error, extraAttrs ...any) {
	status := "success"
	attrs := []any{
		slog.String("status", status),
		slog.Duration("duration", time.Since(r.start)),
	}

	if *err != nil {
		status = "error"
		attrs[0] = slog.String("status", status)
		attrs = append(attrs, slog.Any("error", *err))
	}

	attrs = append(attrs, extraAttrs...)

	r.logger.Debug("Finished "+r.msg, attrs...)
}
