package handler

import (
	"context"
	"go-restfull-api/config"
	"go-restfull-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func EmailExists(ctx context.Context, db *pgxpool.Pool, email string) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1 FROM users WHERE email=$1
        )
    `, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

type User struct {
	ID       string
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DefaultResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}

func HandleSignin(ctx *gin.Context) {
	var loginInput LoginInput

	err := ctx.ShouldBind(&loginInput)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, DefaultResponse{
			Success: false,
			Message: "Unknown error",
		})
		return
	}

	validate := validator.New()

	err = validate.Struct(&loginInput)

	if err != nil {
		errorResponse := util.TransformValidationErrors(err, loginInput)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var user User

	row := config.DB.QueryRow(context.Background(),
		"SELECT * from users WHERE email=$1", loginInput.Email,
	)

	if err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, DefaultResponse{
				Success: false,
				Message: "Invalid email or password",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, DefaultResponse{
			Success: false,
			Message: "Unknown error",
		})
		return
	}

	err = util.CheckPasswordHash(loginInput.Password, user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, DefaultResponse{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	var token string

	token, err = util.GenerateToken(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, DefaultResponse{
			Success: false,
			Message: "Unknown error",
		})
		return
	}

	ctx.JSON(http.StatusOK, DefaultResponse{
		Success: true,
		Message: "Success login",
		Data: map[string]string{
			"token": token,
		},
	})
}

func HandleSignup(ctx *gin.Context) {
	var user User

	err := ctx.ShouldBind(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	validate := validator.New()

	err = validate.Struct(&user)

	if err != nil {
		errorResponse := util.TransformValidationErrors(err, user)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var exists bool

	if exists, err = EmailExists(context.Background(), config.DB, user.Email); exists {
		ctx.JSON(http.StatusConflict, DefaultResponse{
			Success: false,
			Message: "Email already exists",
		})
		return
	}

	var hashedPassword string

	hashedPassword, err = util.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, DefaultResponse{
			Success: false,
			Message: "Unknown error",
		})
		return
	}

	err = config.DB.QueryRow(context.Background(),
		"INSERT INTO users (name, email, password) values ($1, $2, $3) RETURNING id",
		user.Name, user.Email, hashedPassword,
	).Scan(&user.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}

func HandleGetProfile(ctx *gin.Context) {}
