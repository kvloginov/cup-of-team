package handlers

import (
	"net/http"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
)

// Handlers holds dependencies for HTTP handlers
type Handlers struct {
}

// NewHandlers creates a new Handlers instance
func NewHandlers() *Handlers {
	return &Handlers{}
}

// HandleHealth handles GET /health
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	httpServer.SendJSON(w, http.StatusOK, model.HealthResponse{Status: "ok"})
}
