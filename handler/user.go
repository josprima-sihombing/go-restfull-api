package handler

import (
	"context"
	"fmt"
	"go-restfull-api/config"
	"go-restfull-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       string
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func HandleSignin(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "signin endpoint",
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

	// TODO:
	// 2. Check apakah user exist atau ngak by email
	// 3. Hash the password before safe to database
	// 4. Save user to database

	validate := validator.New()

	err = validate.Struct(&user)

	if err != nil {
		errorResponse := util.TransformValidationErrors(err, user)
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	// query user table by email return 0 rows
	// check apakah email exist

	rows, _ := config.DB.Exec(context.Background(),
		"SELECT * from users WHERE email=$1", user.Email,
	)

	err = config.DB.QueryRow(context.Background(),
		"INSERT INTO users (name, email, password) values ($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Password,
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
