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

// CreatePhoto godoc
// @Summary Create photo
// @Description Create photo corresponding to the input
// @Tags photos
// @Accept json
// @Produce json
// @Param models.Photo body models.Photo true "create photo"
// @Success 200 {object} models.Photo
// @Router /photos [post]
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

// GetAllPhoto godoc
// @Summary Get all photo
// @Description Get details of all photo
// @Tags photos
// @Accept json
// @Produce json
// @Success 200 {object} models.Photo
// @Router /photos [get]
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

// GetOnePhoto godoc
// @Summary Get details for a given Id
// @Description Get details of photo corresponding to the input Id
// @Tags photos
// @Accept json
// @Produce json
// @Param Id path int true "ID of the photo"
// @Success 200 {object} models.Photo
// @Router /photos/{Id} [get]
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

// UpdatePhoto godoc
// @Summary Update photo identified by the given Id
// @Description Update the photo corresponding to the input Id
// @Tags photos
// @Accept json
// @Produce json
// @Param Id path int true "ID of the photo to be updated"
// @Success 200 {object} models.Photo
// @Router /photos/{id} [put]
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

// DeletePhoto godoc
// @Summary Delete photo identified by the given Id
// @Description Delete the photo corresponding to the input Id
// @Tags photos
// @Accept json
// @Produce json
// @Param Id path int true "ID of the photo to be deleted"
// @Success 204 "No Content"
// @Router /photos/{id} [delete]
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