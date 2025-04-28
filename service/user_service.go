package service

import (
	"context"
	"go-restfull-api/model"
	"go-restfull-api/repository"
	"go-restfull-api/util"
	"log"
	"net/http"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Signup(c context.Context, user *model.User) (*util.ApiResponse, *model.User) {
	err := util.Validate.Struct(user)

	if err != nil {
		log.Printf("Error: %#v", err)
		errorResponse := util.TransformValidationErrors(err, user)

		return &util.ApiResponse{
			Code: http.StatusBadRequest,
			Data: errorResponse,
		}, nil
	}

	existingUser, err := s.repo.FindByEmail(&c, user.Email)

	if err != nil {
		return util.ServerError, nil
	}

	if existingUser != nil {
		return &util.ApiResponse{
			Code: http.StatusConflict,
		}, nil
	}

	hashedPassword, err := util.HashPassword(user.Password)

	if err != nil {
		return util.ServerError, nil
	}

	newUser := model.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: hashedPassword,
	}

	err, u := s.repo.Save(&c, &newUser)

	if err != nil {
		return util.ServerError, nil
	}

	return nil, u
}

func (s *UserService) Signin(c context.Context, credential *model.UserCredential) (string, *util.ApiResponse) {
	err := util.Validate.Struct(credential)

	if err != nil {
		return "", &util.ApiResponse{
			Code: http.StatusBadRequest,
			Data: util.TransformValidationErrors(err, credential),
		}
	}

	u, err := s.repo.FindByEmail(&c, credential.Email)

	if err != nil {
		return "", util.ServerError
	}

	err = util.CheckPasswordHash(credential.Password, u.Password)

	if err != nil {
		return "", util.ServerError
	}

	token, err := util.GenerateToken(u.ID)

	if err != nil {
		return "", util.ServerError
	}

	return token, nil
}

func (s *UserService) GetProfile(c context.Context, email string) (*model.User, *util.ApiResponse) {
	user, err := s.repo.FindByEmail(&c, email)

	if err != nil {
		return nil, util.ServerError
	}

	return user, nil
}
