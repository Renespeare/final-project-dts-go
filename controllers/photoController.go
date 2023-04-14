package controllers

import (
	// "errors"
	"final-project-dts-go/database"
	"final-project-dts-go/helpers"
	"final-project-dts-go/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePhoto(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		userData = ctx.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(ctx)
		Photo models.Photo
	)

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Photo)
	} else {
		ctx.ShouldBind(&Photo)
	}

	Photo.UserID = userID

	err := db.Debug().Create(&Photo).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&Photo, Photo.ID).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Photo)
}

func GetAllPhoto(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Photos []models.Photo
	)

	err := db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Find(&Photos).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Photos)
}

func GetOnePhoto(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Photo models.Photo
	)

	photoId, _ := strconv.Atoi(ctx.Param("photoId"))

	err := db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&Photo, photoId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Photo)
}

func UpdatePhoto(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(ctx)
		Photo models.Photo
	)

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Photo)
	} else {
		ctx.ShouldBind(&Photo)
	}

	photoId, _ := strconv.Atoi(ctx.Param("photoId"))

	Photo.ID = uint(photoId)

	err := db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Model(&Photo).Where("id = ?", photoId).Updates(models.Photo{Title: Photo.Title, Caption: Photo.Caption, PhotoUrl: Photo.PhotoUrl}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&Photo, photoId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Photo)
}

func DeletePhoto(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Photo models.Photo
	)

	photoId, _ := strconv.Atoi(ctx.Param("photoId"))

	err := db.First(&Photo, "id = ?", photoId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found!",
		})
		return
	}
	err = db.Delete(&Photo).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": "Error deleting photo data",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("Photo with id %v has been successfully deleted", photoId),
	})
}