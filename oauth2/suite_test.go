package oauth2

import (
	"log"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/RichardKnop/go-oauth2-server/api"
	"github.com/RichardKnop/go-oauth2-server/config"
	"github.com/RichardKnop/go-oauth2-server/migrate"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// TestSuite ...
type TestSuite struct {
	suite.Suite
	DB  *gorm.DB
	API *rest.Api
}

// SetupTest creates in-memory test database and starts app
func (suite *TestSuite) SetupTest() {
	if suite.DB == nil {
		db, err := gorm.Open("sqlite3", ":memory:")
		if err != nil {
			log.Fatal(err)
		}
		suite.DB = &db
		migrate.Bootstrap(&db)
		MigrateAll(&db)
	}

	if suite.API == nil {
		suite.API = api.NewAPI(
			api.DevelopmentStack,
			NewRoutes(config.NewConfig(), suite.DB),
		)
	}

	// Insert test client
	clientSecretHash, _ := bcrypt.GenerateFromPassword([]byte("test_client_secret"), 3)
	suite.DB.Create(&Client{
		ClientID: "test_client_id",
		Password: string(clientSecretHash),
	})

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("test_password"), 3)
	// Insert test user
	suite.DB.Create(&User{
		Username:  "test_username",
		Password:  string(passwordHash),
		FirstName: "John",
		LastName:  "Doe",
	})
}

// TearDown truncates all tables
func (suite *TestSuite) TearDown() {
	suite.DB.Exec("DELETE FROM SELECT name FROM sqlite_master WHERE type IS 'table'")
}

// TestTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
