package handler

import (
	"go-restfull-api/model"
	"go-restfull-api/service"
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	baluhap, createdUser := h.service.Signup(ctx, &user)

	if baluhap != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
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

	token, err := h.service.Signin(ctx, &credential)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
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
