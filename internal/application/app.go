package application

import (
	"context"
	"net/http"

	"github.com/bkielbasa/go-web-app/internal/dependency"
	"github.com/gorilla/mux"
)

type App struct {
	httpServer *http.Server
	router     *mux.Router
}

func New(ctx context.Context) *App {
	r := mux.NewRouter()
	healthy := dependency.New()
	r.HandleFunc("/healthyz", healthy.Healthy)
	r.HandleFunc("/readyz", healthy.Ready)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	return &App{
		httpServer: httpServer,
		router:     r,
	}
}

func (app *App) Run() error {
	err := app.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (app *App) Close(ctx context.Context) error {
	return app.httpServer.Shutdown(ctx)
}

type MuxRegister interface {
	MuxRegister(*mux.Router)
}

func (app *App) AddModule(module Module) {
	if m, ok := module.(MuxRegister); ok {
		m.MuxRegister(app.router)
	}
}

type Module interface {
}
