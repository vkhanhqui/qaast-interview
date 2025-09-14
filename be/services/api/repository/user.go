package repository

import (
	"be/pkg/errors"
	"be/pkg/model"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, email, hashed string) (string, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, id, email, name string) (*model.User, error)
	List(ctx context.Context, limit int, cursor string) ([]model.User, error)
}

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) UserRepository {
	return &userRepo{db: pool}
}

func (r *userRepo) Create(ctx context.Context, email, hashed string) (string, error) {
	id := uuid.NewString()
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id,email,password,created_at) VALUES ($1,$2,$3,$4)`,
		id, email, hashed, time.Now().UTC(),
	)
	return id, errors.WithStack(err)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `SELECT id,email,password,name,created_at FROM users WHERE email=$1`, email)
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, errors.WithNotFound(errors.New("User not found"), "")
	}
	return &u, errors.WithStack(err)
}

func (r *userRepo) UpdateUser(ctx context.Context, id, email, name string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(ctx, `
        UPDATE users
        SET
            email = COALESCE(NULLIF($1, ''), email),
            name  = COALESCE(NULLIF($2, ''), name)
        WHERE id = $3
        RETURNING id, email, password, name, created_at
    `, email, name, id).Scan(
		&u.ID,
		&u.Email,
		&u.Password,
		&u.Name,
		&u.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, errors.WithInvalid(errors.New("Email existed"), "")
		}

		if err == pgx.ErrNoRows {
			return nil, errors.WithNotFound(errors.New("User not found"), "")
		}
	}
	return &u, errors.WithStack(err)
}

func (r *userRepo) List(ctx context.Context, limit int, cursor string) ([]model.User, error) {
	query := `
		SELECT id, email, password, name, created_at
		FROM users
	`
	args := []any{}
	if cursor != "" {
		query += ` WHERE id < $1`
		args = append(args, cursor)
	}

	query = fmt.Sprintf("%s ORDER BY id DESC LIMIT $%d", query, len(args)+1)
	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.CreatedAt); err != nil {
			return nil, errors.WithStack(err)
		}
		users = append(users, u)
	}
	return users, errors.WithStack(rows.Err())
}
