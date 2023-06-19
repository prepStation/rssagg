// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Name          string
	Url           string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastFetchedAt sql.NullTime
}

type FeedFollow struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	FeedID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID          uuid.UUID
	Title       string
	Description sql.NullString
	PublishedAt time.Time
	Url         string
	FeedID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ApiKey    string
}
