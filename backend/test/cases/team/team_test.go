package team

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kvloginov/cup-of-team/backend/internal/api/model"
	"github.com/kvloginov/cup-of-team/backend/internal/domain"
	"github.com/kvloginov/cup-of-team/backend/test/env"
	"github.com/stretchr/testify/suite"
)

type TeamTestSuite struct {
	env.BaseSuite
}

func TestTeamSuite(t *testing.T) {
	suite.Run(t, new(TeamTestSuite))
}

// TestTeamFullFlow tests the full flow: create team -> add users -> get team
func (s *TeamTestSuite) TestTeamFullFlow() {
	var teamID string

	// Step 1: Create Team
	s.Run("CreateTeam", func() {
		reqBody := model.CreateTeamRequest{
			Name: "Test Team",
		}
		body, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/api/team", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.Handlers.HandleCreateTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

		var resp model.CreateTeamResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		s.Require().NoError(err)

		s.NotEmpty(resp.ID, "Team ID should not be empty")
		teamID = resp.ID
	})

	// Step 2: Add First User
	s.Run("AddFirstUser", func() {
		reqBody := model.AddToTeamRequest{
			TeamID: teamID,
			User: domain.User{
				ID:                "user1",
				FirstName:         "John",
				Initials:          "JD",
				ParentNames:       []string{"Michael", "Sarah"},
				GrandParentsNames: []string{"Robert", "Mary", "James", "Patricia"},
				Country:           domain.CountryCodeUS,
			},
		}
		body, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/api/team/user", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.Handlers.HandleAddToTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

		var resp model.AddToTeamResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		s.Require().NoError(err)

		s.Equal("user1", resp.User.ID)
		s.Equal("John", resp.User.FirstName)
	})

	// Step 3: Add Second User
	s.Run("AddSecondUser", func() {
		reqBody := model.AddToTeamRequest{
			TeamID: teamID,
			User: domain.User{
				ID:                "user2",
				FirstName:         "Jane",
				Initials:          "JS",
				ParentNames:       []string{"William", "Linda"},
				GrandParentsNames: []string{"George", "Barbara", "Richard", "Susan"},
				Country:           domain.CountryCodeUK,
			},
		}
		body, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/api/team/user", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.Handlers.HandleAddToTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")
	})

	// Step 4: Update First User (test upsert)
	s.Run("UpdateFirstUser", func() {
		reqBody := model.AddToTeamRequest{
			TeamID: teamID,
			User: domain.User{
				ID:                "user1",
				FirstName:         "John Updated",
				Initials:          "JD",
				ParentNames:       []string{"Michael", "Sarah"},
				GrandParentsNames: []string{"Robert", "Mary", "James", "Patricia"},
				Country:           domain.CountryCodeDE,
			},
		}
		body, err := json.Marshal(reqBody)
		s.Require().NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/api/team/user", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		s.Handlers.HandleAddToTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

		var resp model.AddToTeamResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		s.Require().NoError(err)

		s.Equal("user1", resp.User.ID)
		s.Equal("John Updated", resp.User.FirstName)
		s.Equal(domain.CountryCodeDE, resp.User.Country)
	})

	// Step 5: Get Team and Verify All Data
	s.Run("GetTeam", func() {
		url := fmt.Sprintf("/api/team?team_id=%s", teamID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		s.Handlers.HandleGetTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

		var resp model.GetTeamResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		s.Require().NoError(err)

		// Verify team data
		s.Equal(teamID, resp.Team.ID, "Team ID should match")
		s.Equal("Test Team", resp.Team.Name, "Team name should match")

		// Verify users count
		s.Len(resp.Team.Users, 2, "Team should have 2 users")

		// Verify first user (updated)
		var user1, user2 domain.User
		for _, user := range resp.Team.Users {
			if user.ID == "user1" {
				user1 = user
			} else if user.ID == "user2" {
				user2 = user
			}
		}

		s.Equal("user1", user1.ID)
		s.Equal("John Updated", user1.FirstName)
		s.Equal("JD", user1.Initials)
		s.Equal([]string{"Michael", "Sarah"}, user1.ParentNames)
		s.Equal([]string{"Robert", "Mary", "James", "Patricia"}, user1.GrandParentsNames)
		s.Equal(domain.CountryCodeDE, user1.Country)

		// Verify second user
		s.Equal("user2", user2.ID)
		s.Equal("Jane", user2.FirstName)
		s.Equal("JS", user2.Initials)
		s.Equal([]string{"William", "Linda"}, user2.ParentNames)
		s.Equal([]string{"George", "Barbara", "Richard", "Susan"}, user2.GrandParentsNames)
		s.Equal(domain.CountryCodeUK, user2.Country)
	})

	// Step 6: Remove User
	s.Run("RemoveUser", func() {
		url := fmt.Sprintf("/api/team/user?team_id=%s&user_id=%s", teamID, "user2")
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		s.Handlers.HandleRemoveFromTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")
	})

	// Step 7: Verify User Was Removed
	s.Run("VerifyUserRemoved", func() {
		url := fmt.Sprintf("/api/team?team_id=%s", teamID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()

		s.Handlers.HandleGetTeam(w, req)

		s.Equal(http.StatusOK, w.Code, "Expected status 200 OK")

		var resp model.GetTeamResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		s.Require().NoError(err)

		// Verify only 1 user remains
		s.Len(resp.Team.Users, 1, "Team should have only 1 user after removal")
		s.Equal("user1", resp.Team.Users[0].ID, "Remaining user should be user1")
	})
}
