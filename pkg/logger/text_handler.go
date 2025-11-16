package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type colorLogHandler struct {
	baseHandler slog.Handler
	out         io.Writer
	level       slog.Level
	attrs       []slog.Attr
	group       string
}

func newColorLogHandler(w io.Writer, level slog.Level) *colorLogHandler {
	return &colorLogHandler{
		baseHandler: slog.NewTextHandler(w, &slog.HandlerOptions{Level: level}),
		out:         w,
		level:       level,
	}
}

func (h *colorLogHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *colorLogHandler) Handle(_ context.Context, record slog.Record) error {
	if len(h.attrs) > 0 {
		record.AddAttrs(h.attrs...)
	}

	logColor := slogLevelToColor(record.Level)
	levelStr := fmt.Sprintf("%s%s:%s", logColor, strings.ToUpper(record.Level.String()), colorReset)
	timeStr := fmt.Sprintf("[%s]", record.Time.Format("2006-01-02 15:04:05"))

	groupStr := ""
	if h.group != "" {
		groupStr = fmt.Sprintf("[%s]", h.group)
	}

	attrs := make([]string, 0, record.NumAttrs())
	record.Attrs(func(a slog.Attr) bool {
		attrs = append(
			attrs,
			fmt.Sprintf("{%s%s%s=%v}", colorCyan, a.Key, colorReset, a.Value.Any()),
		)

		return true
	})

	attStr := strings.Join(attrs, ",")

	_, err := fmt.Fprintf(
		h.out,
		"%s %s %s %s %s\n",
		timeStr,
		levelStr,
		groupStr,
		record.Message,
		attStr,
	)
	if err != nil {
		return fmt.Errorf("failed write log to output: %w", err)
	}

	return nil
}

func (h *colorLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)

	return &newHandler
}

func (h *colorLogHandler) WithGroup(name string) slog.Handler {
	newHandler := *h

	if h.group == "" {
		newHandler.group = name
	} else {
		newHandler.group = h.group + "." + name
	}

	return &newHandler
}
