package db_provider

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
	config "github.com/spf13/viper"

	db "main/internal/storage/db/generated"
)

type DatabaseProvider struct {
	Conn    *sql.DB
	Queries *db.Queries
}

func NewDatabaseProvider() *DatabaseProvider {
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.GetString("DB_USER"),
		config.GetString("DB_PASSWORD"),
		config.GetString("DB_HOST"),
		config.GetString("DB_PORT"),
		config.GetString("DB_NAME"),
	)

	dbConn, err := sql.Open("postgres", connURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal().Err(err).Msg("failed to ping database")
	}

	dp := &DatabaseProvider{
		Conn:    dbConn,
		Queries: db.New(dbConn),
	}

	return dp
}
