package links

import (
	"github.com/bkielbasa/go-web-app/internal/application"
	"github.com/bkielbasa/go-web-app/links/app"
	"github.com/bkielbasa/go-web-app/links/domain"
	"github.com/bkielbasa/go-web-app/links/infra"
	"github.com/gorilla/mux"
)

func New(shortBaseURL string) application.Module {
	storage := infra.NewInMemoryStorage()
	serv := domain.NewLinksService(shortBaseURL)
	linksApp := app.NewLinks(serv, storage)

	return &linksModule{
		httpHandler: infra.NewHTTPHandler(linksApp),
	}
}

type linksModule struct {
	httpHandler infra.HTTPHandler
}

func (module linksModule) MuxRegister(r *mux.Router) {
	r.HandleFunc("/link", module.httpHandler.Add).Methods("POST")
	r.HandleFunc("/c/{short}", module.httpHandler.Click).Methods("GET")
}
