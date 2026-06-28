package models

import (
	"github.com/google/uuid"
	"time"
)

type Torrent struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	InfoHash    string    `db:"info_hash"`
	Status      int       `db:"status"`
	AuthorID    uuid.UUID `db:"author_id"`
	CategoryID  uuid.UUID `db:"category_id"`
	Downloads   int       `db:"downloads"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
