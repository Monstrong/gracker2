
package dto

import (
	"time"

	"github.com/google/uuid"
)

type Torrent struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	InfoHash    string    `json:"info_hash"`
	AuthorID    uuid.UUID `json:"author_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Status      int32     `json:"status"`
	Downloads   int64     `json:"downloads"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTorrentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InfoHash    string `json:"info_hash"`
	CategoryID  uuid.UUID `json:"category_id"`
}

type CreateTorrentResponse struct {
	Torrent Torrent `json:"torrent"`
}

type GetTorrentResponse struct {
	Torrent Torrent `json:"torrent"`
}

type ListTorrentsResponse struct {
	Torrents []Torrent `json:"torrents"`
}

type UpdateTorrentRequest struct {
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	InfoHash   *string    `json:"info_hash,omitempty"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
}

type UpdateTorrentResponse struct {
	Torrent Torrent `json:"torrent"`
}

type DeleteTorrentResponse struct {
	Deleted bool `json:"deleted"`
}

type UpdateStatusRequest struct {
	NewStatus int32 `json:"new_status"`
}

type UpdateStatusResponse struct {
	Updated bool `json:"updated"`
}

type Category struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type ListCategoryResponse struct {
	Categories []Category `json:"categories"`
}
