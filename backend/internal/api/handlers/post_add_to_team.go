package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	httpServer "github.com/kvloginov/cup-of-team/backend/internal/infra/http"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase"
)

// HandleAddToTeam handles POST /api/team/user
func (h *Handlers) HandleAddToTeam(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req model.AddToTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpServer.SendError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate request
	if req.TeamID == "" {
		httpServer.SendError(w, http.StatusBadRequest, "team_id is required")
		return
	}

	if req.User.ID == "" {
		httpServer.SendError(w, http.StatusBadRequest, "user.id is required")
		return
	}

	if req.User.FirstName == "" {
		httpServer.SendError(w, http.StatusBadRequest, "user.first_name is required")
		return
	}

	log.Printf("[POST /api/team/user] team_id=%s user_id=%s", req.TeamID, req.User.ID)

	// Add user via usecase
	user, err := h.teamUsecase.AddUser(usecase.AddUserParams{
		TeamID: req.TeamID,
		User:   req.User,
	})
	if err != nil {
		httpServer.SendError(w, http.StatusInternalServerError, "Failed to add user to team")
		return
	}

	// Send response
	response := model.AddToTeamResponse{
		User: *user,
	}

	httpServer.SendJSON(w, http.StatusOK, response)
}
