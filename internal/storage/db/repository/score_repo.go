package repository

import (
	"context"

	db_provider "main/internal/storage/db"
	db "main/internal/storage/db/generated"
)

type ScoreRepository interface {
	Create(name string, score int64) error
	Update(name string, score int64) error
	GetAll() ([]*Score, error)
	GetInRange(from, to int64) ([]*Score, error)
	GetScoreByPlayerName(name string) (*Score, error)
}

type Score struct {
	Name  string
	Score int64
	Rank  int64
}

type DefaultScoreRepo struct {
	dbProvider *db_provider.DatabaseProvider
}

func (r *DefaultScoreRepo) GetAll() ([]*Score, error) {
	rows, err := r.dbProvider.Queries.GetAllScores(context.Background())
	if err != nil {
		return nil, err
	}

	var result []*Score
	for _, row := range rows {
		result = append(result, &Score{
			Name:  row.Name.String,
			Score: row.Score,
			Rank:  row.Rank,
		})
	}

	return result, nil
}

func (r *DefaultScoreRepo) GetInRange(from, to int64) ([]*Score, error) {
	rows, err := r.dbProvider.Queries.GetScoresInRange(
		context.Background(),
		db.GetScoresInRangeParams{
			RankFrom: from,
			RankTo:   to,
		},
	)
	if err != nil {
		return nil, err
	}

	var result []*Score
	for _, row := range rows {
		result = append(result, &Score{
			Name:  row.Name.String,
			Score: row.Score,
			Rank:  row.Rank,
		})
	}

	return result, nil
}

func (r *DefaultScoreRepo) GetScoreByPlayerName(name string) (*Score, error) {
	// Player name is unique, so query result will be either empty array or an array with 1 element.
	// This way it's possible to detect if score (and user) exists without getting an error.
	rows, err := r.dbProvider.Queries.GetScoreByPlayerName(
		context.Background(),
		name,
	)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &Score{
		Name:  rows[0].Name.String,
		Score: rows[0].Score,
		Rank:  rows[0].Rank,
	}, nil
}

func (r *DefaultScoreRepo) Create(name string, score int64) error {
	err := r.dbProvider.Queries.CreateScore(
		context.Background(),
		db.CreateScoreParams{
			Score:    score,
			UserName: name,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DefaultScoreRepo) Update(name string, score int64) error {
	err := r.dbProvider.Queries.UpdateScore(
		context.Background(),
		db.UpdateScoreParams{
			Score:    score,
			UserName: name,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
