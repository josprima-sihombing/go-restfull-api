package repository

import (
	"context"
	"errors"
	"go-restfull-api/model"
	"log"

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

	tx, err := r.db.Begin(*c)

	if err != nil {
		log.Printf("Error: %#v", err)
		return err, nil
	}

	err = tx.QueryRow(*c,
		"INSERT INTO users (name, email, password) values ($1, $2, $3) RETURNING id, name, email",
		user.Name, user.Email, user.Password,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email)

	if err != nil {
		log.Printf("Error: %#v", err)
		tx.Rollback(*c)
		return err, nil
	}

	var createdProfile model.Profile

	err = tx.QueryRow(*c,
		"INSERT INTO profiles (user_id) values ($1) RETURNING id, user_id, bio, avatar_url",
		createdUser.ID,
	).Scan(&createdProfile.ID, &createdProfile.UserID, &createdProfile.Bio, &createdProfile.AvatarURL)

	if err != nil {
		log.Printf("Error: %#v", err)
		tx.Rollback(*c)
		return err, nil
	}

	err = tx.Commit(*c)

	if err != nil {
		return err, nil
	}

	return nil, &createdUser
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

func (r *userRepositoryPostgres) UpdateProfile(c *context.Context, id string, profile *model.UpdateProfile) (error, *model.Profile) {
	return nil, nil
}

func (r *userRepositoryPostgres) GetProfile(c *context.Context, id string) (error, *model.ProfileDetail) {
	return nil, nil
}
