package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

const PrivateFieldsMessageString = "NO_INFO"

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) RegisterNewUser(ctx context.Context, user *models.User) (uuid.UUID, error) {
	const op = "repo.Postgres.RegisterNewUser"

	var id uuid.UUID
	query :=
		`
	INSERT INTO users (name, email, password_hash, role, premium, avatar_url, bio)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`
	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Premium,
		user.AvatarURL,
		user.Bio,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrEmailAlreadyExists)
			case "users_name_key":
				return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrNameAlreadyExists)
			default:
				// неизвестный unique-constraint
				return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
			}
		}
		return id, fmt.Errorf("%s: %w", op, err) // 00000000-0000-0000-0000-000000000000
	}

	return id, nil
}

// GetUserPublic gives user PUBLIC PROFILE
func (r *PostgresRepository) GetUserPublic(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const op = "repo.Postgres.GetUserPublic"

	var u models.User
	query := `SELECT id, name, email, role, premium, avatar_url, bio, created_at, updated_at FROM users WHERE id = $1`

	err := r.scanInUwhereX(ctx, query, id, &u)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Email = PrivateFieldsMessageString
	return &u, nil
}

// GetUserPrivate gives user PRIVATE PROFILE
func (r *PostgresRepository) GetUserPrivate(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const op = "repo.Postgres.GetUserPrivate"

	var u models.User
	query := `SELECT id, name, email, role, premium, avatar_url, bio, created_at, updated_at FROM users WHERE id = $1`

	err := r.scanInUwhereX(ctx, query, id, &u)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

// GetUserByEmail gives user PRIVATE PROFILE by email WITH PASSWORD
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "repo.Postgres.GetUserByEmail"

	var u models.User
	query := `SELECT id, name, email, role, premium, avatar_url, bio, created_at, updated_at FROM users WHERE email = $1`
	err := r.scanInUwhereX(ctx, query, email, &u)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = r.db.QueryRow(ctx, `SELECT password_hash FROM users WHERE email = $1`, email).Scan(&u.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

func (r *PostgresRepository) DeleteUser(ctx context.Context, id uuid.UUID) (bool, error) {
	const op = "repo.Postgres.DeleteUser"

	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch tag.RowsAffected() {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default: // тут можно вернуть "предупреждение", о том что удалилось больше одной записи
		return true, nil
	}
}

func (r *PostgresRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role models.Role) (updated bool, err error) {
	const op = "repo.Postgres.UpdateUserRole"

	tag, err := r.db.Exec(ctx, `UPDATE users SET role = $2, 
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $1`, id, role)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch tag.RowsAffected() {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default: // тут можно вернуть "предупреждение", о том что обновилось больше одной записи
		return true, nil
	}
}

func (r *PostgresRepository) UpdateUserPremium(ctx context.Context, id uuid.UUID, premium models.PremiumType) (updated bool, err error) {
	const op = "repo.Postgres.UpdateUserPremium"

	tag, err := r.db.Exec(ctx, `UPDATE users 
	SET premium = $2, 
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $1`, id, premium)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch tag.RowsAffected() {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default: // тут можно вернуть "предупреждение", о том что обновилось больше одной записи
		return true, nil
	}
}

func (r *PostgresRepository) ChangePassword(ctx context.Context, id uuid.UUID, oldPasswordHash string, newPasswordHash string) (updated bool, err error) {
	const op = "repo.Postgres.ChangePassword"
	tag, err := r.db.Exec(ctx, `UPDATE users
	SET password_hash = $2, 
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND password_hash = $3`, id, newPasswordHash, oldPasswordHash)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	switch tag.RowsAffected() {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default: // тут можно вернуть "предупреждение", о том что обновилось больше одной записи
		return true, nil
	}
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, u *models.User) (nu *models.User, err error) {
	const op = "repo.Postgres.UpdateUser"

	// TODO: implement
	return nu, nil
}

// scanInUwhereX scans user in models.User by some parametr. DOES NOT SCAN PASSWORD
func (r *PostgresRepository) scanInUwhereX(ctx context.Context, query string, x any, u *models.User) error {
	err := r.db.QueryRow(ctx, query, x).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Role,
		&u.Premium,
		&u.AvatarURL,
		&u.Bio,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	return err
}


