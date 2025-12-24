package usecase

import "github.com/kvloginov/cup-of-team/backend/internal/domain"

// TeamUsecase defines the interface for team-related business logic
type TeamUsecase interface {
	CreateTeam(params CreateTeamParams) (*CreateTeamResult, error)
	GetTeam(teamID string) (*domain.Team, error)
	AddUser(params AddUserParams) (*domain.User, error)
	RemoveUser(teamID, userID string) error
}

// CreateTeamParams contains parameters for creating a team
type CreateTeamParams struct {
	Name string
}

// CreateTeamResult contains the result of creating a team
type CreateTeamResult struct {
	ID string
}

// AddUserParams contains parameters for adding a user to a team
type AddUserParams struct {
	TeamID string
	User   domain.User
}

