package handlers

import (
	"net/http"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase"
)

// Handlers holds dependencies for HTTP handlers
type Handlers struct {
	teamUsecase usecase.TeamUsecase
}

// NewHandlers creates a new Handlers instance
func NewHandlers(teamUsecase usecase.TeamUsecase) *Handlers {
	return &Handlers{
		teamUsecase: teamUsecase,
	}
}

// HandleHealth handles GET /health
func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	httpServer.SendJSON(w, http.StatusOK, model.HealthResponse{Status: "ok"})
}

type Registerer interface {
	Handle(method, pattern string, handler http.HandlerFunc)
}

func (h *Handlers) RegisterRoutes(server Registerer) {
	server.Handle("POST", "/team", h.HandleCreateTeam)
	server.Handle("GET", "/team", h.HandleGetTeam)
	server.Handle("POST", "/team/user", h.HandleAddToTeam)
	server.Handle("DELETE", "/team/user", h.HandleRemoveFromTeam)

	server.Handle("GET", "/health", h.HandleHealth)
}
