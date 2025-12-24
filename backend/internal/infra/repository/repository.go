package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Team represents a stored team
type Team struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

// User represents a stored user
type User struct {
	ID                string
	TeamID            string
	FirstName         string
	Initials          string
	ParentNames       []string
	GrandParentsNames []string
	Country           string
	CreatedAt         time.Time
}

// Repository handles database operations
type Repository struct {
	db *sql.DB
}

// New creates a new repository
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// ============================================
// TEAM OPERATIONS
// ============================================

// CreateTeam saves a new team to the database
func (r *Repository) CreateTeam(team *Team) error {
	query := `INSERT INTO teams (id, name, created_at)
			  VALUES (?, ?, ?)`

	_, err := r.db.Exec(query,
		team.ID,
		team.Name,
		team.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

// GetTeam retrieves a team by ID
func (r *Repository) GetTeam(id string) (*Team, error) {
	query := `SELECT id, name, created_at
			  FROM teams WHERE id = ?`

	team := &Team{}
	err := r.db.QueryRow(query, id).Scan(
		&team.ID,
		&team.Name,
		&team.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Team not found is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}

// ============================================
// USER OPERATIONS
// ============================================

// CreateUser saves a new user to the database
func (r *Repository) CreateUser(user *User) error {
	parentNamesJSON, err := json.Marshal(user.ParentNames)
	if err != nil {
		return fmt.Errorf("failed to marshal parent_names: %w", err)
	}

	grandParentsNamesJSON, err := json.Marshal(user.GrandParentsNames)
	if err != nil {
		return fmt.Errorf("failed to marshal grandparent_names: %w", err)
	}

	query := `INSERT INTO users (id, team_id, first_name, initials, parent_names, grandparent_names, country, created_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.db.Exec(query,
		user.ID,
		user.TeamID,
		user.FirstName,
		user.Initials,
		string(parentNamesJSON),
		string(grandParentsNamesJSON),
		user.Country,
		user.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (r *Repository) GetUser(id string) (*User, error) {
	query := `SELECT id, team_id, first_name, initials, parent_names, grandparent_names, country, created_at
			  FROM users WHERE id = ?`

	user := &User{}
	var parentNamesJSON string
	var grandParentsNamesJSON string

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.TeamID,
		&user.FirstName,
		&user.Initials,
		&parentNamesJSON,
		&grandParentsNamesJSON,
		&user.Country,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // User not found is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := json.Unmarshal([]byte(parentNamesJSON), &user.ParentNames); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parent_names: %w", err)
	}

	if err := json.Unmarshal([]byte(grandParentsNamesJSON), &user.GrandParentsNames); err != nil {
		return nil, fmt.Errorf("failed to unmarshal grandparent_names: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user in the database
func (r *Repository) UpdateUser(user *User) error {
	parentNamesJSON, err := json.Marshal(user.ParentNames)
	if err != nil {
		return fmt.Errorf("failed to marshal parent_names: %w", err)
	}

	grandParentsNamesJSON, err := json.Marshal(user.GrandParentsNames)
	if err != nil {
		return fmt.Errorf("failed to marshal grandparent_names: %w", err)
	}

	query := `UPDATE users 
			  SET team_id = ?, first_name = ?, initials = ?, parent_names = ?, grandparent_names = ?, country = ?
			  WHERE id = ?`

	_, err = r.db.Exec(query,
		user.TeamID,
		user.FirstName,
		user.Initials,
		string(parentNamesJSON),
		string(grandParentsNamesJSON),
		user.Country,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// DeleteUser removes a user from the database
func (r *Repository) DeleteUser(teamID, userID string) error {
	query := `DELETE FROM users WHERE id = ? AND team_id = ?`

	result, err := r.db.Exec(query, userID, teamID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found or does not belong to team")
	}

	return nil
}

// GetTeamUsers retrieves all users for a team
func (r *Repository) GetTeamUsers(teamID string) ([]User, error) {
	query := `SELECT id, team_id, first_name, initials, parent_names, grandparent_names, country, created_at
			  FROM users WHERE team_id = ? ORDER BY created_at ASC`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var parentNamesJSON string
		var grandParentsNamesJSON string

		if err := rows.Scan(
			&user.ID,
			&user.TeamID,
			&user.FirstName,
			&user.Initials,
			&parentNamesJSON,
			&grandParentsNamesJSON,
			&user.Country,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		if err := json.Unmarshal([]byte(parentNamesJSON), &user.ParentNames); err != nil {
			return nil, fmt.Errorf("failed to unmarshal parent_names: %w", err)
		}

		if err := json.Unmarshal([]byte(grandParentsNamesJSON), &user.GrandParentsNames); err != nil {
			return nil, fmt.Errorf("failed to unmarshal grandparent_names: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}
