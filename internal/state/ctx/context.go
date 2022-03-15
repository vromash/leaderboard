package ctx

import (
	db_provider "main/internal/storage/db"
	"main/internal/storage/db/repository"
)

type AppContext struct {
	DB   *db_provider.DatabaseProvider
	Repo *repository.RepositoryProvider
}

func Setup() *AppContext {
	dbProvider := db_provider.NewDatabaseProvider()
	repoProvider := repository.Setup(dbProvider)

	return &AppContext{
		DB:   dbProvider,
		Repo: repoProvider,
	}
}
