package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"user-photo-api/database"
	"user-photo-api/helpers"
	"user-photo-api/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type photoController struct{}

type PhotoController interface {
	CreatePhotoHandler(c *gin.Context)
	GetPhotoHandler(c *gin.Context)
	UpdatePhotoHandler(c *gin.Context)
	DeletePhotoHandler(c *gin.Context)
}

func NewPhotoController() PhotoController {
	return &photoController{}
}

func (p *photoController) CreatePhotoHandler(c *gin.Context) {
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	tokenUserIdInt, err := strconv.Atoi(tokenUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	var count int64
	database.Database().Model(&models.User{ID: tokenUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	database.Database().Model(&models.Photo{}).Where("user_id = ?", tokenUserIdInt).Count(&count)
	if count > 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("profile photo already exists"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	fileExtension := filepath.Ext(file.Filename)
	if fileExtension != ".jpg" &&
		fileExtension != ".jpeg" &&
		fileExtension != ".png" &&
		fileExtension != ".webp" &&
		fileExtension != ".jfif" {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("photo must be jpg/jpeg/png/webp/jfif"))
		return
	}

	photoFilename := tokenUserIdString + fileExtension
	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, err)
		return
	}
	absolutePhotoFilePath := filepath.Join(rootFolderPath, "uploads", "users", photoFilename)
	protocol := strings.ToLower(strings.Split(c.Request.Proto, "/")[0])
	urlPhoto := protocol + "://" + c.Request.Host + "/uploads/users/" + photoFilename

	err = c.SaveUploadedFile(file, absolutePhotoFilePath)
	if err != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, err)
		return
	}

	photo := new(models.Photo)
	photo.PhotoUrl = urlPhoto
	photo.UserId = tokenUserIdInt
	if err := database.Database().Create(&photo).Error; err != nil {
		errRemove := os.Remove(absolutePhotoFilePath)
		if errRemove != nil {
			helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
			return
		}
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}

	photoResponse := new(models.PhotoResponse)
	copier.Copy(photoResponse, photo)

	helpers.SendResponseSuccess(c, photoResponse)
}

func (p *photoController) GetPhotoHandler(c *gin.Context) {
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	tokenUserIdInt, err := strconv.Atoi(tokenUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}

	var count int64
	database.Database().Model(&models.User{ID: tokenUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	photo := new(models.Photo)
	if err := database.Database().
		Where("user_id = ?", tokenUserIdInt).
		Limit(1).Find(photo).Error; err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	if photo.ID == 0 {
		helpers.SendResponseSuccess(c, "user has no photos")
		return
	}

	photoResponse := new(models.PhotoResponse)
	copier.Copy(photoResponse, photo)
	helpers.SendResponseSuccess(c, photoResponse)
}

func (p *photoController) UpdatePhotoHandler(c *gin.Context) {
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	tokenUserIdInt, err := strconv.Atoi(tokenUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	var count int64
	database.Database().Model(&models.User{ID: tokenUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	paramPhotoIdString := c.Param("photoId")
	paramPhotoIdInt, err := strconv.Atoi(paramPhotoIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}

	paramPhoto := new(models.Photo)
	if err := database.Database().
		Limit(1).Find(paramPhoto, paramPhotoIdInt).Error; err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	if paramPhoto.ID == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("photo not found"))
		return
	}

	if paramPhoto.UserId != tokenUserIdInt {
		helpers.SendResponseError(c, http.StatusUnauthorized, errors.New("every user can only edit own photo"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	fileExtension := filepath.Ext(file.Filename)
	if fileExtension != ".jpg" &&
		fileExtension != ".jpeg" &&
		fileExtension != ".png" &&
		fileExtension != ".webp" &&
		fileExtension != ".jfif" {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("photo must be jpg/jpeg/png/webp/jfif"))
		return
	}

	newPhotoFilename := tokenUserIdString + fileExtension
	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, err)
		return
	}
	newPhotoFilePath := filepath.Join(rootFolderPath, "uploads", "users", newPhotoFilename)
	protocol := strings.ToLower(strings.Split(c.Request.Proto, "/")[0])
	newUrlPhoto := protocol + "://" + c.Request.Host + "/uploads/users/" + newPhotoFilename

	oldPhotoUrl := paramPhoto.PhotoUrl
	oldPhotoFilename := path.Base(oldPhotoUrl)
	oldPhotoFilePath := filepath.Join(rootFolderPath, "uploads", "users", oldPhotoFilename)

	tempPhotoFilePath := filepath.Join(rootFolderPath, "uploads", "temp", oldPhotoFilename)

	err = CopyFile(oldPhotoFilePath, tempPhotoFilePath)
	if err != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, err)
		return
	}

	if newPhotoFilePath == oldPhotoFilePath {
		err = c.SaveUploadedFile(file, newPhotoFilePath)
		if err != nil {
			errRemove := os.Remove(tempPhotoFilePath)
			if errRemove != nil {
				helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
				return
			}

			helpers.SendResponseError(c, http.StatusInternalServerError, err)
			return
		}

		errRemove := os.Remove(tempPhotoFilePath)
		if errRemove != nil {
			helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
			return
		}

		photoResponse := new(models.PhotoResponse)
		copier.Copy(photoResponse, paramPhoto)
		helpers.SendResponseSuccess(c, photoResponse)
	} else {
		err = c.SaveUploadedFile(file, newPhotoFilePath)
		if err != nil {
			helpers.SendResponseError(c, http.StatusInternalServerError, err)
			return
		}

		if err := database.Database().
			Model(paramPhoto).Update("photo_url", newUrlPhoto).Error; err != nil {
			errRemove := os.Remove(newPhotoFilePath)
			if errRemove != nil {
				helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
				return
			}

			errRemove = os.Remove(tempPhotoFilePath)
			if errRemove != nil {
				helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
				return
			}

			helpers.SendResponseError(c, http.StatusBadRequest, err)
			return
		}

		errRemove := os.Remove(oldPhotoFilePath)
		if errRemove != nil {
			helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
			return
		}

		errRemove = os.Remove(tempPhotoFilePath)
		if errRemove != nil {
			helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
			return
		}

		photoResponse := new(models.PhotoResponse)
		copier.Copy(photoResponse, paramPhoto)
		helpers.SendResponseSuccess(c, photoResponse)
	}
}

func (p *photoController) DeletePhotoHandler(c *gin.Context) {
	tokenUserIdAny, exists := c.Get("tokenUserIdString")
	if !exists {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("ID not found in context"))
		return
	}
	tokenUserIdString := fmt.Sprintf("%v", tokenUserIdAny)
	tokenUserIdInt, err := strconv.Atoi(tokenUserIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	var count int64
	database.Database().Model(&models.User{ID: tokenUserIdInt}).Count(&count)
	if count == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("user not found"))
		return
	}

	paramPhotoIdString := c.Param("photoId")
	paramPhotoIdInt, err := strconv.Atoi(paramPhotoIdString)
	if err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}

	paramPhoto := new(models.Photo)
	if err := database.Database().
		Limit(1).Find(paramPhoto, paramPhotoIdInt).Error; err != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, err)
		return
	}
	if paramPhoto.ID == 0 {
		helpers.SendResponseError(c, http.StatusBadRequest, errors.New("photo not found"))
		return
	}
	if paramPhoto.UserId != tokenUserIdInt {
		helpers.SendResponseError(c, http.StatusUnauthorized, errors.New("every user can only delete own photo"))
		return
	}

	oldPhotoUrl := paramPhoto.PhotoUrl
	oldPhotoFilename := path.Base(oldPhotoUrl)
	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, err)
		return
	}
	oldPhotoFilePath := filepath.Join(rootFolderPath, "uploads", "users", oldPhotoFilename)

	res := database.Database().
		Delete(&models.Photo{}, paramPhoto.ID)
	if res.Error != nil {
		helpers.SendResponseError(c, http.StatusBadRequest, res.Error)
		return
	}

	errRemove := os.Remove(oldPhotoFilePath)
	if errRemove != nil {
		helpers.SendResponseError(c, http.StatusInternalServerError, errRemove)
		return
	}

	helpers.SendResponseSuccess(c, "")
}

func CopyFile(srcPath, dstPath string) error {
	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Sync the destination file to ensure it's written to disk
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
