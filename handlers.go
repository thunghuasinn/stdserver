package stdserver

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Ctx struct {
	fiber.Ctx
}

func (app *App) Use(handlers ...fiber.Handler) *App {
	for _, h := range handlers {
		app.fibre.Use(h)
	}
	return app
}

func (app *App) Get(path string, handler fiber.Handler) *App {
	app.fibre.Get(path, handler)
	return app
}

func (app *App) Post(path string, handler fiber.Handler) *App {
	app.fibre.Post(path, handler)
	return app
}

func errorHandler(c *fiber.Ctx, err error) error {
	// TODO error log?
	code := fiber.StatusInternalServerError
	prefix := ""
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	} else {
		prefix = "[UNHANDLED] "
	}
	return c.Status(code).JSON(fiber.Map{
		"message": prefix + err.Error(),
		"object":  fmt.Sprintf("%+v", err),
	})
}
