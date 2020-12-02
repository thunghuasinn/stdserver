package stdserver

import (
	"context"
	"net/http"
)
import "github.com/gofiber/adaptor/v2"
import "github.com/gofiber/fiber/v2"

func (app *App) UseHttpHandler(h http.Handler) {
	app.Use(func(c *fiber.Ctx) error {
		return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", c.Locals("dbUser"))))
		})(c)
	})
}
