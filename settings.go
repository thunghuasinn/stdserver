package stdserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Settings struct {
	fiber.Config
	Limiter limiter.Config
	Name         string
	CSRFPath     string
	LoginHandler LoginFunc
	KeyTableDir  string
}
