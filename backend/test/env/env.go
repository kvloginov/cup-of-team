package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kvloginov/cup-of-team/backend/internal/api/handlers"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/db"
	"github.com/kvloginov/cup-of-team/backend/internal/infra/repository"
	"github.com/kvloginov/cup-of-team/backend/internal/usecase/team"
	"github.com/stretchr/testify/suite"
)

type BaseSuite struct {
	suite.Suite
	DB       *db.DB
	Repo     *repository.Repository
	Usecase  *team.Usecase
	Handlers *handlers.Handlers
	dbPath   string
}

func (s *BaseSuite) SetupTest() {
}

func TestBaseSuite(t *testing.T) {
	suite.Run(t, new(BaseSuite))
}

func (s *BaseSuite) SetupSuite() {
	// Create temporary database file for tests
	s.dbPath = filepath.Join(os.TempDir(), "cup-of-team-test.db")

	database, err := db.New(s.dbPath)
	s.Require().NoError(err, "Failed to initialize test database")
	s.DB = database

	s.Repo = repository.New(database.DB)
	s.Usecase = team.NewUsecase(s.Repo)
	s.Handlers = handlers.NewHandlers(s.Usecase)
}

func (s *BaseSuite) TearDownSuite() {
	// Close database
	if s.DB != nil {
		s.DB.Close()
	}

	// Remove test database file
	if s.dbPath != "" {
		os.Remove(s.dbPath)
	}
}
