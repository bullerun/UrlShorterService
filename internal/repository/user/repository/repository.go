package user

import (
	"UrlShorterService/internal/entity"
	userStorage "UrlShorterService/internal/repository/user"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	createTable(pool)
	return &Repository{pool: pool}
}

func createTable(pool *pgxpool.Pool) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR (50) UNIQUE NOT NULL,
    email VARCHAR (300) UNIQUE NOT NULL,
    password VARCHAR (255) NOT NULL,
    salt varchar(50)
);`
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := pool.Exec(ctx, createTableQuery)
	if err != nil {
		log.Fatalf("error creating user table: %s", err)
	}
}

func (r *Repository) AddUser(ctx context.Context, user *entity.User) (int64, error) {
	const op = "user.repository.AddUser"
	row := r.pool.QueryRow(ctx, "Insert into users (username, email, password,salt) values ($1, $2, $3, $4) Returning id",
		user.Name,
		user.Email,
		user.Password,
		user.Salt)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, userStorage.ErrUserAlreadyExist
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "user.repository.GetUser"

	row := r.pool.QueryRow(ctx, "select users.id, users.email, users.username, users.salt, users.password from users where users.email = $1", email)
	var user entity.User
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.Salt, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, userStorage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil

}
