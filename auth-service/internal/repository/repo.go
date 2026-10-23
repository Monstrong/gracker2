package repository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/monstrong/gracker2/auth-service/internal/models"
	"github.com/monstrong/gracker2/pkg/errs"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresRepository struct {
	db *pgxpool.Pool
}

const PrivateFieldsMessageString = "NO_INFO"
const userColumns = `id, name, email, role, premium, avatar_url, bio, created_at, updated_at`

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
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
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
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, userColumns)

	err := r.scanInUwhereX(ctx, query, id, &u)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	u.Email = PrivateFieldsMessageString
	return &u, nil
}

// GetUserPrivate gives user PRIVATE PROFILE
func (r *PostgresRepository) GetUserPrivate(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const op = "repo.Postgres.GetUserPrivate"

	var u models.User
	query := fmt.Sprintf(`SELECT %s FROM users WHERE id = $1`, userColumns)

	err := r.scanInUwhereX(ctx, query, id, &u)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &u, nil
}

// GetUserByEmail gives user PRIVATE PROFILE by email WITH PASSWORD
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "repo.Postgres.GetUserByEmail"

	var u models.User
	query := fmt.Sprintf(`SELECT %s, password_hash FROM users WHERE email = $1`, userColumns)
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Role, &u.Premium,
		&u.AvatarURL, &u.Bio, &u.CreatedAt, &u.UpdatedAt, &u.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
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
	return tag.RowsAffected() == 1, nil
}

func (r *PostgresRepository) UpdateUserRole(ctx context.Context, id uuid.UUID, role models.Role) (updated bool, err error) {
	const op = "repo.Postgres.UpdateUserRole"

	tag, err := r.db.Exec(ctx, `UPDATE users SET role = $2, 
	updated_at = CURRENT_TIMESTAMP
	WHERE id = $1`, id, role)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return tag.RowsAffected() == 1, nil
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
	return tag.RowsAffected() == 1, nil
}

// GetPasswordHash возвращает старый хеш пароля
func (r *PostgresRepository) GetPasswordHash(ctx context.Context, id uuid.UUID) (string, error) {
	const op = "repo.Postgres.GetPasswordHash"

	var hash string
	err := r.db.QueryRow(ctx, `SELECT password_hash FROM users WHERE id = $1`, id).Scan(&hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return hash, nil
}

func (r *PostgresRepository) ChangePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) (updated bool, err error) {
	const op = "repo.Postgres.ChangePassword"

	tag, err := r.db.Exec(ctx,
		`UPDATE users
            SET password_hash = $2,
                updated_at = CURRENT_TIMESTAMP
          WHERE id = $1`,
		id, newPasswordHash,
	)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return tag.RowsAffected() == 1, nil
}
func (r *PostgresRepository) UpdateUser(ctx context.Context, u *models.UserPatch, id uuid.UUID) (*models.User, error) {
	const op = "repo.Postgres.UpdateUser"

	b := sq.Update("users")

	if u.Name == nil && u.Email == nil && u.AvatarURL == nil && u.Bio == nil {
		return nil, fmt.Errorf("%s: %w", op, errs.ErrUserPatchEmpty)
	}
	if u.Name != nil {
		b = b.Set("name", *u.Name)
	}
	if u.Email != nil {
		b = b.Set("email", *u.Email)
	}
	if u.AvatarURL != nil {
		b = b.Set("avatar_url", *u.AvatarURL)
	}
	if u.Bio != nil {
		b = b.Set("bio", *u.Bio)
	}
	b = b.Set("updated_at", sq.Expr("CURRENT_TIMESTAMP"))

	b.Where(sq.Eq{"id": id})
	b = b.Suffix(fmt.Sprintf(`RETURNING %s`, userColumns))

	query, args, err := b.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	var nu models.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&nu.ID,
		&nu.Name,
		&nu.Email,
		&nu.Role,
		&nu.Premium,
		&nu.AvatarURL,
		&nu.Bio,
		&nu.CreatedAt,
		&nu.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, errs.ErrNotFound)
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, fmt.Errorf("%s: %w", op, errs.ErrEmailAlreadyExists)
			case "users_name_key":
				return nil, fmt.Errorf("%s: %w", op, errs.ErrNameAlreadyExists)
			default:
				// неизвестный unique-constraint
				return nil, fmt.Errorf("%s: %w", op, errs.ErrAlreadyExists)
			}
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &nu, nil

}

func (r *PostgresRepository) CreateRefreshToken(ctx context.Context, id uuid.UUID, token string) (bool, error) {
	const op = "repo.Postgres.CreateRefreshToken"
	tag, err := r.db.Exec(ctx, `INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2)`, id, token)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return tag.RowsAffected() == 1, nil
}

func (r *PostgresRepository) RotateRefreshToken(ctx context.Context, id uuid.UUID, token string) error {

}

func (r *PostgresRepository) RevokeRefreshTokens(ctx context.Context, id uuid.UUID, token string) (bool, error) {

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
