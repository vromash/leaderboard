package ctx

import db_provider "main/internal/storage/db"

type AppContext struct {
	DB *db_provider.DatabaseProvider
}

func Setup() *AppContext {
	dbProvider := db_provider.NewDatabaseProvider()

	return &AppContext{
		DB: dbProvider,
	}
}
