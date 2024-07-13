package app

// InitWorkers initializes worker.
func (a *App) initWorkers() []worker {
	workers := []worker{
		serveHttp,
		serveBroker,
	}

	return workers
}
