package repository

import (
	db_provider "main/internal/storage/db"
)

type RepositoryProvider struct {
	User  UserRepository
	Score ScoreRepository
}

func Setup(
	dbProvider *db_provider.DatabaseProvider,
) *RepositoryProvider {
	return &RepositoryProvider{
		User:  &DefaultUserRepo{dbProvider},
		Score: &DefaultScoreRepo{dbProvider},
	}
}
