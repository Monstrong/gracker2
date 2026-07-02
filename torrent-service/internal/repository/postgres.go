package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/pkg/logger"
)

type PostgresRepository struct {
	db     *pgxpool.Pool
	l logger.Logger
}

func NewPostgresRepository(db *pgxpool.Pool, logger logger.Logger) *PostgresRepository {
	return &PostgresRepository{
		db:     db,
		l: logger,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, t *models.Torrent) (*models.Torrent, error) {
	const op = "repo.postgres.create"
	query := `
	INSERT INTO torrents (name, description, info_hash, author_id, category_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING (id, name, description, info_hash, status, author_id, category_id, downloads, created_at, updated_at)
	`
	err := r.db.QueryRow(ctx, query, 
		t.Name, 
		t.Description, 
		t.InfoHash, 
		t.AuthorID, 
		t.CategoryID,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.InfoHash,
		&t.Status,
		&t.AuthorID,
		&t.CategoryID,
		&t.Downloads,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return t, nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*models.Torrent, error) {
	// TODO: add filters
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error) {

}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {

}
