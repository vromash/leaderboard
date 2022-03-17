// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createScore = `-- name: CreateScore :exec
INSERT INTO "score" ("score", "user_id")
VALUES ($1, (SELECT "id" FROM "user" WHERE "name" = $2::varchar))
`

type CreateScoreParams struct {
	Score    int64
	UserName string
}

func (q *Queries) CreateScore(ctx context.Context, arg CreateScoreParams) error {
	_, err := q.db.ExecContext(ctx, createScore, arg.Score, arg.UserName)
	return err
}

const createUser = `-- name: CreateUser :exec

INSERT INTO "user" ("name")
VALUES ($1)
`

//- User ---
func (q *Queries) CreateUser(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, createUser, name)
	return err
}

const getAllScores = `-- name: GetAllScores :many

SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub
`

type GetAllScoresRow struct {
	Name  sql.NullString
	Score int64
	Rank  int64
}

//- Score ---
func (q *Queries) GetAllScores(ctx context.Context) ([]GetAllScoresRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllScores)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllScoresRow
	for rows.Next() {
		var i GetAllScoresRow
		if err := rows.Scan(&i.Name, &i.Score, &i.Rank); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecordNumber = `-- name: GetRecordNumber :one
SELECT COUNT(id)
FROM "score"
`

func (q *Queries) GetRecordNumber(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, getRecordNumber)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getRecordNumberInTimeRange = `-- name: GetRecordNumberInTimeRange :one
SELECT COUNT(id)
FROM "score"
WHERE "updated_at" > $1::timestamp
`

func (q *Queries) GetRecordNumberInTimeRange(ctx context.Context, updatedAt time.Time) (int64, error) {
	row := q.db.QueryRowContext(ctx, getRecordNumberInTimeRange, updatedAt)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getScoreByPlayerName = `-- name: GetScoreByPlayerName :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT user_id,
                score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
     ) AS sub
WHERE sub.user_id = (SELECT "id" FROM "user" WHERE "name" = $1::varchar)
`

type GetScoreByPlayerNameRow struct {
	Name  sql.NullString
	Score int64
	Rank  int64
}

func (q *Queries) GetScoreByPlayerName(ctx context.Context, userName string) ([]GetScoreByPlayerNameRow, error) {
	rows, err := q.db.QueryContext(ctx, getScoreByPlayerName, userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetScoreByPlayerNameRow
	for rows.Next() {
		var i GetScoreByPlayerNameRow
		if err := rows.Scan(&i.Name, &i.Score, &i.Rank); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScoreByPlayerNameInTimeRange = `-- name: GetScoreByPlayerNameInTimeRange :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT user_id,
                score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
         WHERE "updated_at" > $1
     ) AS sub
WHERE sub.user_id = (SELECT "id" FROM "user" WHERE "name" = $2::varchar)
`

type GetScoreByPlayerNameInTimeRangeParams struct {
	DateTo   time.Time
	UserName string
}

type GetScoreByPlayerNameInTimeRangeRow struct {
	Name  sql.NullString
	Score int64
	Rank  int64
}

func (q *Queries) GetScoreByPlayerNameInTimeRange(ctx context.Context, arg GetScoreByPlayerNameInTimeRangeParams) ([]GetScoreByPlayerNameInTimeRangeRow, error) {
	rows, err := q.db.QueryContext(ctx, getScoreByPlayerNameInTimeRange, arg.DateTo, arg.UserName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetScoreByPlayerNameInTimeRangeRow
	for rows.Next() {
		var i GetScoreByPlayerNameInTimeRangeRow
		if err := rows.Scan(&i.Name, &i.Score, &i.Rank); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScoresInRange = `-- name: GetScoresInRange :many
SELECT sub.name, sub.score, sub.rank
FROM (
         SELECT score,
                u.name,
                ROW_NUMBER() OVER (ORDER BY "score" DESC) AS rank
         FROM "score"
                  LEFT JOIN "user" u ON u.id = score.user_id
         WHERE "updated_at" > $1
     ) AS sub
WHERE sub.rank BETWEEN $2::bigint AND $3::bigint
`

type GetScoresInRangeParams struct {
	DateTo   time.Time
	RankFrom int64
	RankTo   int64
}

type GetScoresInRangeRow struct {
	Name  sql.NullString
	Score int64
	Rank  int64
}

func (q *Queries) GetScoresInRange(ctx context.Context, arg GetScoresInRangeParams) ([]GetScoresInRangeRow, error) {
	rows, err := q.db.QueryContext(ctx, getScoresInRange, arg.DateTo, arg.RankFrom, arg.RankTo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetScoresInRangeRow
	for rows.Next() {
		var i GetScoresInRangeRow
		if err := rows.Scan(&i.Name, &i.Score, &i.Rank); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateScore = `-- name: UpdateScore :exec
UPDATE "score"
SET "score" = $1
WHERE "user_id" = (SELECT "id" FROM "user" WHERE "name" = $2::varchar)
`

type UpdateScoreParams struct {
	Score    int64
	UserName string
}

func (q *Queries) UpdateScore(ctx context.Context, arg UpdateScoreParams) error {
	_, err := q.db.ExecContext(ctx, updateScore, arg.Score, arg.UserName)
	return err
}
