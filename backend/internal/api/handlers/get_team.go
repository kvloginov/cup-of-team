package handlers

import (
	"log"
	"net/http"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
)

// HandleGetTeam handles GET /api/team
func (h *Handlers) HandleGetTeam(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	teamID := r.URL.Query().Get("team_id")
	if teamID == "" {
		httpServer.SendError(w, http.StatusBadRequest, "team_id parameter is required")
		return
	}

	log.Printf("[GET /api/team] team_id=%s", teamID)

	// Get team via usecase
	team, err := h.teamUsecase.GetTeam(teamID)
	if err != nil {
		httpServer.SendError(w, http.StatusNotFound, "Team not found")
		return
	}

	// Send response
	response := model.GetTeamResponse{
		Team: *team,
	}

	httpServer.SendJSON(w, http.StatusOK, response)
}
