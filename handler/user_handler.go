package handler

import (
	"encoding/json"
	"errors"
	"go-restfull-api/model"
	"go-restfull-api/service"
	"go-restfull-api/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) HandleSignup(ctx *gin.Context) {
	var user model.User

	err := ctx.ShouldBind(&user)

	if err != nil {
		log.Printf("Error: %#v\n", err)

		var syntaxErr *json.SyntaxError

		if errors.As(err, &syntaxErr) {
			ctx.JSON(http.StatusBadRequest, &util.ApiResponse{
				Code: http.StatusBadRequest,
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, util.ServerError)

		return
	}

	signupError, createdUser := h.service.Signup(ctx, &user)

	if signupError != nil {
		ctx.JSON(signupError.Code, signupError)

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": createdUser,
	})
}

func (h *UserHandler) HandleSignin(ctx *gin.Context) {
	var credential model.UserCredential

	err := ctx.ShouldBind(&credential)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	token, signinError := h.service.Signin(ctx, &credential)

	if signinError != nil {
		ctx.JSON(http.StatusBadRequest, signinError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *UserHandler) HandleGetProfile(ctx *gin.Context) {
	user, err := h.service.GetProfile(ctx, "")

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
