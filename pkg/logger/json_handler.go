package logger

import (
	"context"
	"io"
	"log/slog"
	"time"
)

type jsonLogHandler struct {
	baseHandler slog.Handler
	level       slog.Level
}

func newJSONLogHandler(w io.Writer, level slog.Level) *jsonLogHandler {
	return &jsonLogHandler{
		baseHandler: slog.NewJSONHandler(
			w,
			&slog.HandlerOptions{
				Level: level,
				ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
					if attr.Key == slog.TimeKey {
						attr.Value = slog.StringValue(attr.Value.Time().Format(time.DateTime))
						return attr
					}

					if attr.Value.Kind() == slog.KindDuration {
						attr.Value = slog.StringValue(attr.Value.Duration().String())
					}

					return attr
				},
			},
		),
		level: level,
	}
}

func (h *jsonLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *jsonLogHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.baseHandler.Handle(ctx, record)
}

func (h *jsonLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &jsonLogHandler{
		baseHandler: h.baseHandler.WithAttrs(attrs),
		level:       h.level,
	}
}

func (h *jsonLogHandler) WithGroup(name string) slog.Handler {
	return &jsonLogHandler{
		baseHandler: h.baseHandler.WithGroup(name),
		level:       h.level,
	}
}
