package stdserver

import "net/http"
import "github.com/gofiber/adaptor/v2"

func (app *App) UseHttpHandler(h http.Handler) {
	app.Use(adaptor.HTTPHandler(h))
}
