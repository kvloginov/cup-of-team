package team

import (
	"fmt"
	"time"

	"github.com/kvloginov/cup-of-team/backend/internal/domain"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/repository"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase"
)

// Usecase handles team-related business logic
type Usecase struct {
	repo *repository.Repository
}

// NewUsecase creates a new team Usecase instance
func NewUsecase(repo *repository.Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

// CreateTeam creates a new team
func (u *Usecase) CreateTeam(params usecase.CreateTeamParams) (*usecase.CreateTeamResult, error) {
	// Generate unique team ID
	teamID := fmt.Sprintf("team_%d", time.Now().UnixNano())

	// Create team in database
	team := &repository.Team{
		ID:        teamID,
		Name:      params.Name,
		CreatedAt: time.Now(),
	}

	if err := u.repo.CreateTeam(team); err != nil {
		return nil, fmt.Errorf("failed to create team: %w", err)
	}

	return &usecase.CreateTeamResult{
		ID: teamID,
	}, nil
}

// GetTeam retrieves a team with all its users
func (u *Usecase) GetTeam(teamID string) (*domain.Team, error) {
	// Get team
	team, err := u.repo.GetTeam(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	if team == nil {
		return nil, fmt.Errorf("team not found")
	}

	// Get team users
	users, err := u.repo.GetTeamUsers(teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}

	// Convert repository users to domain users
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = domain.User{
			ID:                user.ID,
			FirstName:         user.FirstName,
			Initials:          user.Initials,
			ParentNames:       user.ParentNames,
			GrandParentsNames: user.GrandParentsNames,
			Country:           user.Country,
		}
	}

	return &domain.Team{
		ID:    team.ID,
		Name:  team.Name,
		Users: domainUsers,
	}, nil
}

// AddUser adds or updates a user in a team
func (u *Usecase) AddUser(params usecase.AddUserParams) (*domain.User, error) {
	// Verify team exists
	team, err := u.repo.GetTeam(params.TeamID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify team: %w", err)
	}

	if team == nil {
		return nil, fmt.Errorf("team not found")
	}

	// Check if user already exists
	existingUser, err := u.repo.GetUser(params.User.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		// Update existing user
		user := &repository.User{
			ID:                params.User.ID,
			TeamID:            params.TeamID,
			FirstName:         params.User.FirstName,
			Initials:          params.User.Initials,
			ParentNames:       params.User.ParentNames,
			GrandParentsNames: params.User.GrandParentsNames,
			Country:           params.User.Country,
			CreatedAt:         existingUser.CreatedAt, // Preserve original creation time
		}

		if err := u.repo.UpdateUser(user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}

		return &domain.User{
			ID:                user.ID,
			FirstName:         user.FirstName,
			Initials:          user.Initials,
			ParentNames:       user.ParentNames,
			GrandParentsNames: user.GrandParentsNames,
			Country:           user.Country,
		}, nil
	}

	// Create new user
	user := &repository.User{
		ID:                params.User.ID,
		TeamID:            params.TeamID,
		FirstName:         params.User.FirstName,
		Initials:          params.User.Initials,
		ParentNames:       params.User.ParentNames,
		GrandParentsNames: params.User.GrandParentsNames,
		Country:           params.User.Country,
		CreatedAt:         time.Now(),
	}

	if err := u.repo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &domain.User{
		ID:                user.ID,
		FirstName:         user.FirstName,
		Initials:          user.Initials,
		ParentNames:       user.ParentNames,
		GrandParentsNames: user.GrandParentsNames,
		Country:           user.Country,
	}, nil
}

// RemoveUser removes a user from a team
func (u *Usecase) RemoveUser(teamID, userID string) error {
	// Verify team exists
	team, err := u.repo.GetTeam(teamID)
	if err != nil {
		return fmt.Errorf("failed to verify team: %w", err)
	}

	if team == nil {
		return fmt.Errorf("team not found")
	}

	// Delete user
	if err := u.repo.DeleteUser(teamID, userID); err != nil {
		return fmt.Errorf("failed to remove user: %w", err)
	}

	return nil
}

