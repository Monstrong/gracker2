package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, t *models.Torrent) (*models.Torrent, error) {
	const op = "repo.postgres.create"
	query := `
	INSERT INTO torrents (name, description, info_hash, author_id, category_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, name, description, info_hash, status, author_id, category_id, downloads, created_at, updated_at
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
	const op = "repo.postgres.list"
	var ts []*models.Torrent
	query := `SELECT id, name, description, info_hash, status, author_id, category_id, downloads, created_at, updated_at FROM torrents 
	ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Torrent
		err := rows.Scan(
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
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ts = append(ts, &t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return ts, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Torrent, error) {
	const op = "repo.postgres.getbyid"

	var t models.Torrent
	query := `SELECT id, name, description, info_hash, status, author_id, category_id, downloads, created_at, updated_at FROM torrents WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, models.ErrNotFound)
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return &t, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	const op = "repo.postgres.delete"

	tag, err := r.db.Exec(ctx, `DELETE FROM torrents WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch tag.RowsAffected() {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return true, fmt.Errorf("%s: more than 1 rows deleted", op)
	}
}
