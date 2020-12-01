package stdserver

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	defaultAppName = "KStdServer"
)

type App struct {
	fibre        *fiber.App
	settings     *Settings
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
	app := &App{
		fibre:    fiber.New(s.Config),
		settings: s,
		logger:   s.Logger,
	}
	app.loggingEntry = app.logger.WithField("app", s.Name)
	app.settings.Logger = app.loggingEntry
	app.initMiddlewares()
	return app
}

func (app *App) Log(module string) *logrus.Entry {
	return app.loggingEntry.WithField("module", module)
}

func (app *App) Start(addr string) error {
	log := app.Log("core/main")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Info("starting server...")
		if err := app.fibre.Listen(addr); err != nil {
			errs <- err
		} else {
			done <- true
		}
	}()

	select {
	case sig := <-sigs:
		log.Infof("signal %s received, shutting down...", sig)
		if err := app.fibre.Shutdown(); err != nil {
			log.WithError(err).Error()
			return err
		}
		log.Info("server stopped gracefully")
	case <-done:
		log.Info("server stopped gracefully")
	case err := <-errs:
		log.WithError(err).Error()
		return err
	}
	log.Info("goodbye!")
	return nil
}
