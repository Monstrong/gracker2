package repository

import (
	"context"

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

func (r *CategoriesRepository) IsExists(ctx context.Context, categoryID uuid.UUID) (bool, error) {
	const op = "repo.categories.isexists"
}

func (r *CategoriesRepository) List(ctx context.Context, limit, offset int) ([]*models.Category, error) {
	const op = "repo.categories.list"
}
