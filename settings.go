package stdserver

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/rs/zerolog"
)

type Settings struct {
	Fiber        fiber.Config
	Limiter      limiter.Config
	Name         string
	CSRFPath     string
	LoginPath    string
	LoginHandler LoginFunc
	KeyTableDir  string
	Logger       *zerolog.Logger
	Context      context.Context
	SkipAuth     func(ctx *fiber.Ctx) bool
	AllowOrigins string

	// LogLevel sets the default accepted level for logging. It is ignored if Settings.Logger is provided.
	LogLevel zerolog.Level

	ColorfulLogging bool
}
