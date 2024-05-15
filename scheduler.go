package stdserver

import (
	"context"
	"time"
)

type Task func(ctx context.Context)

func (app *App) StartChild(task Task) {
	app.lock.RLock()
	defer app.lock.RUnlock()
	app.children = append(app.children, task)
	if app.running {
		go task(app.ctx)
	}
}

func (app *App) StartChildInterval(interval time.Duration, task, teardown Task) {
	app.StartChild(func(ctx context.Context) {
		t := time.NewTicker(interval)
	outer:
		for {
			select {
			case <-ctx.Done():
				break outer
			case <-t.C:
				task(ctx)
			}
		}
		t.Stop()
		if teardown != nil {
			// since main ctx is already closed
			teardown(context.TODO())
		}
	})
}

func (app *App) StartChildDelayed(delay time.Duration, task, teardown Task) {
	app.StartChild(func(ctx context.Context) {
		t := time.NewTimer(delay)
		select {
		case <-ctx.Done():
		case <-t.C:
			task(ctx)
		}
		t.Stop()
		if teardown != nil {
			// since main ctx is already closed
			teardown(context.TODO())
		}
	})
}
