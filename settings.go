package stdserver

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/limiter"
)

type Settings struct {
	fiber.Settings
	limiter.Config
	Name         string
	CSRFPath     string
	LoginHandler LoginFunc
	KeyTableDir  string
}
