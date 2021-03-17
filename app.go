package stdserver

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	defaultAppName = "KStdServer"
)

type App struct {
	fibre        *fiber.App
	settings     *Settings
	rootCtx      context.Context
	ctx          context.Context
	cancelCtx    context.CancelFunc
	logger       logrus.FieldLogger
	loggingEntry logrus.FieldLogger
}

func New(settings ...*Settings) *App {
	s := &Settings{}
	s.ETag = false
	s.Prefork = true
	if len(settings) > 0 {
		s = settings[0]
	}
	if len(s.Name) == 0 {
		s.Name = defaultAppName
	}
	if len(s.ServerHeader) == 0 {
		s.ServerHeader = s.Name
	}
	if s.ErrorHandler == nil {
		s.ErrorHandler = errorHandler
	}
	if s.Logger == nil {
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})
		s.Logger = l
	}
	if s.IdleTimeout == 0 {
		s.IdleTimeout = 10 * time.Second
	}
	if s.Context == nil {
		s.Context = context.TODO()
	}
	app := &App{
		fibre:    fiber.New(s.Config),
		settings: s,
		logger:   s.Logger,
		rootCtx:  s.Context,
	}
	app.loggingEntry = app.logger.WithField("app", s.Name)
	app.settings.Logger = app.loggingEntry
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("app", app)
		c.Locals("config", app.settings)
		return c.Next()
	})
	app.initMiddlewares()
	return app
}

func (app *App) Log(module string) *logrus.Entry {
	return app.loggingEntry.WithField("module", module)
}

func (app *App) Start(addr string) error {
	return app.WrapStart(func() error {
		return app.fibre.Listen(addr)
	})
}

func (app *App) StartTLS(addr, certFile, keyFile string) error {
	return app.WrapStart(func() error {
		return app.fibre.ListenTLS(addr, certFile, keyFile)
	})
}

func (app *App) WrapStart(startFunc func() error) error {
	log := app.Log("core/main")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	app.ctx, app.cancelCtx = context.WithCancel(app.rootCtx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Info("starting server...")
		if err := startFunc(); err != nil {
			errs <- err
		} else {
			done <- true
		}
		app.cancelCtx()
	}()

	select {
	case sig := <-sigs:
		log.Infof("signal %s received, shutting down...", sig)
		app.cancelCtx()
		if err := app.fibre.Shutdown(); err != nil {
			log.WithError(err).Error()
			return err
		}
		log.Info("server stopped gracefully")
	case <-done:
		app.cancelCtx()
		log.Info("server stopped gracefully")
	case err := <-errs:
		app.cancelCtx()
		log.WithError(err).Error()
		return err
	}
	log.Info("goodbye!")
	return nil
}

func (app *App) Router(prefix string, handlers ...fiber.Handler) fiber.Router {
	return app.fibre.Group(prefix, handlers...)
}
