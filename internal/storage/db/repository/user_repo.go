package repository

import (
	"context"

	db_provider "main/internal/storage/db"
)

type UserRepository interface {
	Create(name string) error
}

type User struct {
	Name string
}

type DefaultUserRepo struct {
	dbProvider *db_provider.DatabaseProvider
}

func (r *DefaultUserRepo) Create(name string) error {
	err := r.dbProvider.Queries.CreateUser(
		context.Background(),
		name,
	)
	if err != nil {
		return err
	}

	return nil
}
