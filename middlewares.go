package stdserver

import (
	"github.com/gofiber/cors"
	"github.com/gofiber/csrf"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/limiter"
)

var (
	defaultCSRFPath      = "/csrf"
	defaultLimiterConfig = limiter.Config{
		Timeout: 1,
		Max:     5,
	}
)

func (app *App) initMiddlewares() {
	s := app.settings
	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(cors.New())
	app.Use(middleware.Compress())
	// app.Use(middleware.Favicon())
	if s.Config.Timeout == 0 {
		s.Config.Timeout = defaultLimiterConfig.Timeout
	}
	if s.Config.Max == 0 {
		s.Config.Max = defaultLimiterConfig.Max
	}
	app.Use(limiter.New(s.Config))
	app.Use(csrf.New())
	if len(s.CSRFPath) == 0 {
		s.CSRFPath = defaultCSRFPath
	}
	app.Get(s.CSRFPath, csrfHandler)
}

func csrfHandler(c *fiber.Ctx) {
	_ = c.JSON(fiber.Map{"csrf": c.Locals("csrf")})
}
