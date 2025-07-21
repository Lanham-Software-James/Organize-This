// Package tests is where all of out unit tests are described.
package mocks

import (
	"willowsuite-vault/infra/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockDB returns a gorm DB and sqlMock object to use during unit testing.
func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logger.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		logger.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}
