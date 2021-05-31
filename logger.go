package stdserver

import (
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func newLoggerMiddleware(cfg *Settings) fiber.Handler {
	log := cfg.Logger.With().Str("module", "core").Str("subModule", "fiber").Logger()
	pid := os.Getpid()
	var once sync.Once
	var errHandler fiber.ErrorHandler

	return func(c *fiber.Ctx) error {
		once.Do(func() {
			errHandler = c.App().Config().ErrorHandler
		})

		start := time.Now()
		chainErr := c.Next()
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}
		stop := time.Now()

		logCtx := log.With().
			Int("pid", pid).
			Str("remote", c.IP()).
			Int("status", c.Response().StatusCode()).
			Int("length", len(c.Response().Body())).
			Int64("latency", stop.Sub(start).Milliseconds()).
			Str("userAgent", c.Get(fiber.HeaderUserAgent)).
			Str("method", c.Method()).
			Str("url", c.OriginalURL())
		if user := c.Locals("user"); user != nil {
			logCtx = logCtx.Str("user", user.(string))
		}
		logger := logCtx.Logger()
		logger.Err(chainErr).Msg("request")

		return nil
	}
}
