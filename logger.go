package stdserver

import (
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func newLoggerMiddleware(cfg *Settings) fiber.Handler {
	log := cfg.Logger.WithField("module", "fiber")
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

		fields := make(logrus.Fields)
		fields["pid"] = pid
		fields["remote"] = c.IP()
		if user := c.Locals("user"); user != nil {
			fields["user"] = user
		}
		fields["status"] = c.Response().StatusCode()
		fields["length"] = len(c.Response().Body())
		fields["latency"] = stop.Sub(start).Milliseconds()
		fields["userAgent"] = c.Get(fiber.HeaderUserAgent)

		if chainErr != nil {
			log.WithFields(fields).WithError(chainErr).Errorf("%s %s", c.Method(), c.OriginalURL())
		} else {
			log.WithFields(fields).Infof("%s %s", c.Method(), c.OriginalURL())
		}

		return nil
	}
}
