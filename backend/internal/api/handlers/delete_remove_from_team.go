package handlers

import (
	"log"
	"net/http"

	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
)

// HandleRemoveFromTeam handles DELETE /api/team/user
func (h *Handlers) HandleRemoveFromTeam(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	teamID := r.URL.Query().Get("team_id")
	userID := r.URL.Query().Get("user_id")

	if teamID == "" {
		httpServer.SendError(w, http.StatusBadRequest, "team_id parameter is required")
		return
	}

	if userID == "" {
		httpServer.SendError(w, http.StatusBadRequest, "user_id parameter is required")
		return
	}

	log.Printf("[DELETE /api/team/user] team_id=%s user_id=%s", teamID, userID)

	// Remove user via usecase
	if err := h.teamUsecase.RemoveUser(teamID, userID); err != nil {
		httpServer.SendError(w, http.StatusInternalServerError, "Failed to remove user from team")
		return
	}

	// Send empty response
	httpServer.SendJSON(w, http.StatusOK, struct{}{})
}
