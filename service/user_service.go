package service

import (
	"context"
	"errors"
	"go-restfull-api/model"
	"go-restfull-api/repository"
	"go-restfull-api/util"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Signup(c context.Context, user *model.User) (any, *model.User) {
	err := util.Validate.Struct(user)

	if err != nil {
		errorResponse := util.TransformValidationErrors(err, user)
		return &errorResponse, nil
	}

	hashedPassword, err := util.HashPassword(user.Password)

	if err != nil {
		return errors.New("error"), nil
	}

	newUser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	err, u := s.repo.Save(&c, &newUser)

	if err != nil {
		return errors.New("Unexpected error"), nil
	}

	return nil, u
}

func (s *UserService) Signin(c context.Context, credential *model.UserCredential) (any, error) {
	err := util.Validate.Struct(credential)

	if err != nil {
		return nil, errors.New("Unexpected error")
	}

	u, err := s.repo.FindByEmail(&c, credential.Email)

	if err != nil {
		return nil, errors.New("Unexpected error")
	}

	err = util.CheckPasswordHash(credential.Password, u.Password)

	if err != nil {
		return nil, errors.New("Unexpected error")
	}

	token, err := util.GenerateToken(u.ID)

	if err != nil {
		return nil, errors.New("Unexpected error")
	}

	return token, nil
}

func (s *UserService) GetProfile(c context.Context, email string) (*model.User, error) {
	user, err := s.repo.FindByEmail(&c, email)

	if err != nil {
		return nil, errors.New("Unexpected error")
	}

	return user, nil
}
