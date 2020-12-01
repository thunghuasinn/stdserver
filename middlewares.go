package stdserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"time"
)

var (
	defaultCSRFPath      = "/csrf"
	defaultLimiterConfig = limiter.Config{
		Expiration: 1 * time.Second,
		Max:        5,
	}
)

func (app *App) initMiddlewares() {
	s := app.settings
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{AllowCredentials: true}))
	app.Use(compress.New())
	// app.Use(middleware.Favicon())
	if s.Limiter.Expiration == 0 {
		s.Limiter.Expiration = defaultLimiterConfig.Expiration
	}
	if s.Limiter.Max == 0 {
		s.Limiter.Max = defaultLimiterConfig.Max
	}
	app.Use(limiter.New(s.Limiter))
	app.Use(csrf.New(csrf.Config{ContextKey: "csrf"}))
	if len(s.CSRFPath) == 0 {
		s.CSRFPath = defaultCSRFPath
	}
	app.Get(s.CSRFPath, csrfHandler)
}

func csrfHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"data": fiber.Map{"csrf": c.Locals("csrf")}})
}
