package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase"
)

// HandleCreateTeam handles POST /api/team
func (h *Handlers) HandleCreateTeam(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req model.CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpServer.SendError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate request
	if req.Name == "" {
		httpServer.SendError(w, http.StatusBadRequest, "Team name is required")
		return
	}

	log.Printf("[POST /api/team] name=%s", req.Name)

	// Create team via usecase
	result, err := h.teamUsecase.CreateTeam(usecase.CreateTeamParams{
		Name: req.Name,
	})
	if err != nil {
		httpServer.SendError(w, http.StatusInternalServerError, "Failed to create team")
		return
	}

	// Send response
	response := model.CreateTeamResponse{
		ID: result.ID,
	}

	httpServer.SendJSON(w, http.StatusOK, response)
}
