package infra

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/bkielbasa/go-web-app/links/app"
	"github.com/bkielbasa/go-web-app/links/domain"
	"github.com/gorilla/mux"
)

type AddLinkRequest struct {
	Target string
	Tags   []string
}

type AddLinkResponse struct {
	ShortURL string
}

type HTTPHandler struct {
	links app.Link
}

func NewHTTPHandler(linksApp app.Link) HTTPHandler {
	return HTTPHandler{links: linksApp}
}

func (h HTTPHandler) Add(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := AddLinkRequest{}
	if err = json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL, err := h.links.Add(r.Context(), req.Target, req.Tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, _ := json.Marshal(AddLinkResponse{
		ShortURL: shortURL,
	})
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(respBody)
}

func (h HTTPHandler) Click(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL := vars["short"]

	targetLink, err := h.links.ClickLink(r.Context(), shortURL)

	if errors.Is(err, domain.ErrLinkNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, targetLink, http.StatusMovedPermanently)
}
