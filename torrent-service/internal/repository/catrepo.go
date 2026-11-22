package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/torrent-service/internal/models"
)

type CategoriesRepository struct {
	db *pgxpool.Pool
}

func NewCategoriesRepository(db *pgxpool.Pool) *CategoriesRepository {
	return &CategoriesRepository{
		db: db,
	}
}

func (r *CategoriesRepository) IsExists(
	ctx context.Context,
	categoryID uuid.UUID,
) (bool, error) {
	const op = "repo.categories.isexists"

	query := `
        SELECT EXISTS (
            SELECT 1
            FROM categories
            WHERE id = $1
        )
    `

	var exists bool
	if err := r.db.QueryRow(ctx, query, categoryID).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (r *CategoriesRepository) List(ctx context.Context) ([]*models.Category, error) {
	const op = "repo.categories.list"

	query := `
        SELECT id, name, parent_id
        FROM categories
        ORDER BY name ASC
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	categories := make([]*models.Category, 0)

	for rows.Next() {
		category := new(models.Category)

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID,
		); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows: %w", op, err)
	}

	return categories, nil
}
