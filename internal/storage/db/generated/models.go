// Code generated by sqlc. DO NOT EDIT.

package db

import ()

type Score struct {
	ID     int64
	Score  int64
	UserID int64
}

type User struct {
	ID   int64
	Name string
}