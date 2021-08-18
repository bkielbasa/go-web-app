package infra

import "net/http"

type Dependency interface {
	Healthy() bool
	Ready() bool
}

type DependencyManager struct {
	dependencies []Dependency
}

func NewHealthy() *DependencyManager {
	return &DependencyManager{
		dependencies: []Dependency{},
	}
}

func (h *DependencyManager) Healthy(w http.ResponseWriter, r *http.Request) {
	for _, s := range h.dependencies {
		if !s.Healthy() {
			http.Error(w, "Unhealthy", http.StatusInternalServerError)
		}
	}
}

func (h *DependencyManager) Ready(w http.ResponseWriter, r *http.Request) {
	for _, s := range h.dependencies {
		if !s.Ready() {
			http.Error(w, "Not ready", http.StatusInternalServerError)
		}
	}
}
