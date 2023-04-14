package controllers

import (
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

// CreateSocialMedia godoc
// @Summary Create social media
// @Description Create social media corresponding to the input
// @Tags social_media
// @Accept json
// @Produce json
// @Param models.SocialMedia body models.SocialMedia true "create social media"
// @Success 200 {object} models.SocialMedia
// @Router /social_media [post]
func CreateSocialMedia(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		userData = ctx.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(ctx)
		SocialMedia models.SocialMedia
	)

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&SocialMedia)
	} else {
		ctx.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&SocialMedia, SocialMedia.ID).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, SocialMedia)
}

// GetAllSocialMedia godoc
// @Summary Get all social media
// @Description Get details of all social media
// @Tags social_media
// @Accept json
// @Produce json
// @Success 200 {object} models.SocialMedia
// @Router /social_media [get]
func GetAllSocialMedia(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		SocialMedia []models.SocialMedia
	)

	err := db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Find(&SocialMedia).Error
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SocialMedia)
}

// GetOneSocialMedia godoc
// @Summary Get details for a given Id
// @Description Get details of social media corresponding to the input Id
// @Tags social_media
// @Accept json
// @Produce json
// @Param Id path int true "ID of the social media"
// @Success 200 {object} models.SocialMedia
// @Router /social_media/{Id} [get]
func GetOneSocialMedia(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		SocialMedia models.SocialMedia
	)

	socialMediaId, _ := strconv.Atoi(ctx.Param("socialMediaId"))

	err := db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&SocialMedia, socialMediaId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SocialMedia)
}

// UpdateSocialMedia godoc
// @Summary Update social media identified by the given Id
// @Description Update the social media corresponding to the input Id
// @Tags social_media
// @Accept json
// @Produce json
// @Param Id path int true "ID of the social media to be updated"
// @Success 200 {object} models.SocialMedia
// @Router /social_media/{id} [put]
func UpdateSocialMedia(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(ctx)
		SocialMedia models.SocialMedia
	)

	if contentType == appJSON {
		ctx.ShouldBindJSON(&SocialMedia)
	} else {
		ctx.ShouldBind(&SocialMedia)
	}

	socialMediaId, _ := strconv.Atoi(ctx.Param("socialMediaId"))

	SocialMedia.ID = uint(socialMediaId)

	err := db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Model(&SocialMedia).Where("id = ?", socialMediaId).Updates(models.SocialMedia{Name: SocialMedia.Name, SocialMediaUrl: SocialMedia.SocialMediaUrl}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).First(&SocialMedia, socialMediaId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SocialMedia)
}

// DeleteSocialMedia godoc
// @Summary Delete social media identified by the given Id
// @Description Delete the social media corresponding to the input Id
// @Tags social_media
// @Accept json
// @Produce json
// @Param Id path int true "ID of the social media to be deleted"
// @Success 204 "No Content"
// @Router /social_media/{id} [delete]
func DeleteSocialMedia(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		SocialMedia models.SocialMedia
	)

	socialMediaId, _ := strconv.Atoi(ctx.Param("socialMediaId"))

	err := db.First(&SocialMedia, "id = ?", socialMediaId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found!",
		})
		return
	}
	err = db.Delete(&SocialMedia).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": "Error deleting social media data",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("Social media with id %v has been successfully deleted", socialMediaId),
	})
}