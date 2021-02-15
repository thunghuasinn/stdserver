package stdserver

import (
	"context"
	"time"
)

func (app *App) Context() context.Context {
	return app.ctx
}
func (app *App) ContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(app.ctx)
}

func (app *App) ContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(app.ctx, timeout)
}

func (app *App) ContextWithDeadline(deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(app.ctx, deadline)
}

func (app *App) ContextWithValue(key interface{}, value interface{}) context.Context {
	return context.WithValue(app.ctx, key, value)
}
