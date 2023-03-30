package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"user-photo-api/database"
	"user-photo-api/helpers"
	"user-photo-api/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type userController struct{}

type UserController interface {
	RegisterUserHandler(c *gin.Context)
	UpdateUserHandler(c *gin.Context)
	DeleteUserHandler(c *gin.Context)
}

func NewUserController() UserController {
	return &userController{}
}

func (u *userController) RegisterUserHandler(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	req.Trim()
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	user := new(models.User)
	copier.Copy(user, req)
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	user.Password = hashedPassword
	if err := database.Database().Create(&user).Error; err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	userRegisterResponse := new(models.UserRegisterResponse)
	copier.Copy(userRegisterResponse, user)
	helpers.SendResponseSuccess(c, userRegisterResponse)
}

func (u *userController) UpdateUserHandler(c *gin.Context) {
	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	req.Trim()
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	paramUserIdString := c.Param("userId")
	if paramUserIdString != tokenUserIdString {
		helpers.SendResponseError(c, http.StatusUnauthorized, errors.New("every user can only edit own profile"))
		return
	}
	paramUserIdInt, err := strconv.Atoi(paramUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	var count int64
	database.Database().Model(&models.User{ID: paramUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}
	user := new(models.User)
	copier.Copy(user, req)
	if err := database.Database().
		Model(&models.User{ID: paramUserIdInt}).Updates(user).Error; err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	userResponse := new(models.UserResponse)
	copier.Copy(userResponse, user)
	helpers.SendResponseSuccess(c, userResponse)
}

func (u *userController) DeleteUserHandler(c *gin.Context) {
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	paramUserIdString := c.Param("userId")
	if paramUserIdString != tokenUserIdString {
		helpers.SendResponseError(c, http.StatusUnauthorized, errors.New("every user can only delete own profile"))
		return
	}
	paramUserIdInt, err := strconv.Atoi(paramUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	var count int64
	database.Database().Model(&models.User{ID: paramUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}
	res := database.Database().
		Delete(&models.User{}, paramUserIdInt)
	if res.Error != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, res.Error)
		return
	}
	helpers.SendResponseSuccess(c, "")
}
