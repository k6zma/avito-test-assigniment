package middlewares

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NewHTTPLogger(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()

		attrs := []any{
			slog.String("method", c.Method()),
			slog.String("url", c.Path()),
			slog.Int("http_code", status),
			slog.String("ip", c.IP()),
			slog.Duration("latency", latency),
		}

		sanitize := func(s string) string {
			s = strings.ReplaceAll(s, "\\n", ": ")
			s = strings.ReplaceAll(s, "\n", ": ")
			s = strings.ReplaceAll(s, " : ", ": ")

			return s
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", sanitize(err.Error())))
		} else if status >= 400 {
			if cause, ok := c.Locals("error_cause").(error); ok && cause != nil {
				attrs = append(attrs, slog.String("error_cause", sanitize(cause.Error())))
			}

			body := string(c.Response().Body())
			if body != "" {
				attrs = append(attrs, slog.String("error", sanitize(body)))
			}
		}

		switch {
		case status >= 500:
			log.Error("http request", attrs...)
		case status >= 400:
			log.Warn("http request", attrs...)
		default:
			log.Info("http request", attrs...)
		}

		return err
	}
}
