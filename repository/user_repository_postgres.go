package repository

import (
	"context"
	"errors"
	"go-restfull-api/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepositoryPostgres struct {
	db *pgxpool.Pool
}

func NewUserRepositoryPostgres(db *pgxpool.Pool) UserRepository {
	return &userRepositoryPostgres{
		db: db,
	}
}

func (r *userRepositoryPostgres) Save(c *context.Context, user *model.User) (error, *model.User) {
	var createdUser model.User

	err := r.db.QueryRow(*c,
		"INSERT INTO users (name, email, password) values ($1, $2, $3) RETURNING id, name, email",
		user.Name, user.Email, user.Password,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email)

	return err, &createdUser
}

func (r *userRepositoryPostgres) FindByEmail(c *context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(*c,
		"SELECT id, name, email FROM users WHERE email=$1", email,
	).Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
