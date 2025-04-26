package repository

import (
	"context"
	"go-restfull-api/model"
)

type UserRepository interface {
	Save(c *context.Context, user *model.User) (error, *model.User)
	FindByEmail(c *context.Context, email string) (*model.User, error)
}
