package app

import (
	"context"
	"sync"
)

// runWorkers run workers.
func (a *App) runWorkers(ctx context.Context) {
	workers := a.initWorkers()

	wg := new(sync.WaitGroup)
	wg.Add(len(workers))

	for _, work := range workers {
		go func(ctx context.Context, work func(context.Context, *App), t *App) {
			work(ctx, t)
			wg.Done()
		}(ctx, work, a)
	}

	wg.Wait()
}
