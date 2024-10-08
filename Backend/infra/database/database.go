package database

import (
	"log"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	database *gorm.DB
	once     sync.Once
	err      error
)

// DbConnection create database connection
func DbConnection(masterDSN, replicaDSN string) error {
	var db = database
	var err error

	once.Do(func() {
		logMode := viper.GetBool("DB_LOG_MODE")
		debug := viper.GetBool("DEBUG")

		loglevel := logger.Silent
		if logMode {
			loglevel = logger.Info
		}

		db, err = gorm.Open(postgres.Open(masterDSN), &gorm.Config{
			Logger: logger.Default.LogMode(loglevel),
		})
		if !debug {
			db.Use(dbresolver.Register(dbresolver.Config{
				Replicas: []gorm.Dialector{
					postgres.Open(replicaDSN),
				},
				Policy: dbresolver.RandomPolicy{},
			}))
		}
		if err != nil {
			log.Fatalf("Db connection error: %s", err)
		}
		database = db
	})

	return err
}

// GetDB connection
func GetDB() *gorm.DB {
	return database
}
