package stdserver

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	fiber.Config
	Limiter      limiter.Config
	Name         string
	CSRFPath     string
	LoginPath    string
	LoginHandler LoginFunc
	KeyTableDir  string
	Logger       logrus.FieldLogger
	Context      context.Context
	SkipAuth     func(ctx *fiber.Ctx) bool
}
