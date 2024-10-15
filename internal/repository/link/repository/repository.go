package link

import (
	linkStorage "UrlShorterService/internal/repository/link"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, url, alias string, id int64) error {
	const op = "link.repository.Save"
	_, err := r.db.Exec(ctx, "INSERT INTO links (url, alias, users_id) VALUES ($1, $2, $3)", url, alias, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return linkStorage.ErrAliasAlreadyExist
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
