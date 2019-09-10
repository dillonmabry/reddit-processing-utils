package datamanager

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/dillonmabry/reddit-processing-utils/src/logging"
)

var db *sql.DB

// InitDB initialize db
func InitDB(dataSourceName string) {
	logger := logging.NewLogger()
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		logger.Fatal("Could not open database connection for persistence")
	}

	if err = db.Ping(); err != nil {
		logger.Fatal("Database connectivity issue, ping unsuccessful")
	}
	logger.Info("Connected to database")
}
