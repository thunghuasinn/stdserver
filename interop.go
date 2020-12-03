package stdserver

import (
	"net/http"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
)

func (app *App) UseHttpHandler(h http.Handler) {
	app.Use(adaptor.HTTPHandler(h))
}

func MuxSubRouter(parent fiber.Router, prefix string) *mux.Router {
	child := mux.NewRouter()
	parent.Group(prefix, adaptor.HTTPHandler(child))
	return child
}
