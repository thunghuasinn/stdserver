package stdserver

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	defaultAppName = "KStdServer"
)

type App struct {
	fibre     *fiber.App
	settings  *Settings
	rootCtx   context.Context
	ctx       context.Context
	cancelCtx context.CancelFunc
	logger    zerolog.Logger
	running   bool
	lock      sync.RWMutex
	children  []Task
}

func New(settings ...*Settings) *App {
	s := &Settings{}
	s.Fiber.ETag = false
	s.Fiber.Prefork = true
	if len(settings) > 0 {
		s = settings[0]
	}
	if len(s.Name) == 0 {
		s.Name = defaultAppName
	}
	if len(s.Fiber.ServerHeader) == 0 {
		s.Fiber.ServerHeader = s.Name
	}
	if s.Fiber.ErrorHandler == nil {
		s.Fiber.ErrorHandler = errorHandler
	}
	if s.Logger == nil {
		var w io.Writer = os.Stderr
		if s.ColorfulLogging {
			w = zerolog.ConsoleWriter{Out: w}
		}
		l := zerolog.New(w).Level(s.LogLevel)
		s.Logger = &l
	}
	if s.Fiber.IdleTimeout == 0 {
		s.Fiber.IdleTimeout = 10 * time.Second
	}
	if s.Context == nil {
		s.Context = context.TODO()
	}
	app := &App{
		fibre:    fiber.New(s.Fiber),
		settings: s,
		logger:   s.Logger.With().Str("app", s.Name).Logger(),
		rootCtx:  s.Context,
	}
	s.Logger = &app.logger
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("app", app)
		c.Locals("config", app.settings)
		return c.Next()
	})
	app.initMiddlewares()
	return app
}

func (app *App) Fiber() *fiber.App {
	return app.fibre
}

func (app *App) Log(module string, subModule ...string) *zerolog.Logger {
	l := app.logger.With().Str("module", module)
	if len(subModule) > 0 {
		l = l.Str("subModule", subModule[0])
	}
	nl := l.Logger()
	return &nl
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

func (app *App) IsRunning() bool {
	app.lock.RLock()
	defer app.lock.RUnlock()
	return app.running
}

func (app *App) setRunning(state bool) {
	app.lock.Lock()
	app.running = state
	app.lock.Unlock()
}

func (app *App) WrapStart(startFunc func() error) error {
	app.lock.Lock()
	log := app.Log("core", "main")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	app.ctx, app.cancelCtx = context.WithCancel(app.rootCtx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Info().Msg("starting server...")
		if err := startFunc(); err != nil {
			errs <- err
		} else {
			done <- true
		}
		app.setRunning(false)
		app.cancelCtx()
	}()
	for _, child := range app.children {
		go child(app.ctx)
	}
	app.running = true
	app.lock.Unlock()

	select {
	case sig := <-sigs:
		log.Info().Msgf("signal %s received, shutting down...", sig)
		app.cancelCtx()
		app.setRunning(false)
		if err := app.fibre.Shutdown(); err != nil {
			log.Err(err).Msg("while shutting down")
			return err
		}
		log.Info().Msg("server stopped gracefully")
	case <-done:
		app.cancelCtx()
		app.setRunning(false)
		log.Info().Msg("server stopped gracefully")
	case err := <-errs:
		app.cancelCtx()
		log.Err(err).Msg("while running")
		app.setRunning(false)
		return err
	}
	log.Info().Msg("goodbye!")
	return nil
}

func (app *App) Router(prefix string, handlers ...fiber.Handler) fiber.Router {
	return app.fibre.Group(prefix, handlers...)
}
