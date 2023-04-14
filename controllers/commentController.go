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

// CreateComment godoc
// @Summary Create comment
// @Description Create comment corresponding to the input
// @Tags comments
// @Accept json
// @Produce json
// @Param models.Comment body models.Comment true "create comment"
// @Success 200 {object} models.Comment
// @Router /comments [post]
func CreateComment(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		userData = ctx.MustGet("userData").(jwt.MapClaims)
		contentType = helpers.GetContentType(ctx)
		Comment models.Comment
	)

	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Comment)
	} else {
		ctx.ShouldBind(&Comment)
	}

	photoId, ok := ctx.GetQuery("photo_id")
	if !ok || photoId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message":  "Please input photo_id in query param",
		})
		return
	} 

	Photo := models.Photo{}
	err := db.Debug().First(&Photo, photoId).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
			"message":  "Photo not found, can't be commented",
		})
		return
	}

	Comment.UserID = userID
	Comment.PhotoID = Photo.ID

	err = db.Debug().Create(&Comment).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Preload("Photo",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Title", "Caption", "PhotoUrl")
	}).First(&Comment, Comment.ID).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, Comment)
}

// GetAllBook godoc
// @Summary Get all comment
// @Description Get details of all comment with given photo id
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {object} models.Comment
// @Router /comments [get]
func GetAllComment(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Comments []models.Comment
	)

	photoId, ok := ctx.GetQuery("photo_id")
	if !ok || photoId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message":  "Please input photo_id in query param",
		})
		return
	}

	Photo := models.Photo{}
	err := db.Debug().First(&Photo, photoId).Error
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Not Found",
			"message":  "Photo not found, can't see comments",
		})
		return
	}

	err = db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Preload("Photo",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Title", "Caption", "PhotoUrl")
	}).Where("photo_id = ?", photoId).Find(&Comments).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Comments)
}

// GetOneComment godoc
// @Summary Get details for a given Id
// @Description Get details of comment corresponding to the input Id
// @Tags comments
// @Accept json
// @Produce json
// @Param Id path int true "ID of the comment"
// @Success 200 {object} models.Comment
// @Router /comments/{Id} [get]
func GetOneComment(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Comment models.Comment
	)

	commentId, _ := strconv.Atoi(ctx.Param("commentId"))

	err := db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Preload("Photo",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Title", "Caption", "PhotoUrl")
	}).First(&Comment, commentId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Comment)
}

// UpdateComment godoc
// @Summary Update comment identified by the given Id
// @Description Update the comment corresponding to the input Id
// @Tags comments
// @Accept json
// @Produce json
// @Param Id path int true "ID of the comment to be updated"
// @Success 200 {object} models.Comment
// @Router /comments/{id} [put]
func UpdateComment(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(ctx)
		Comment models.Comment
	)

	if contentType == appJSON {
		ctx.ShouldBindJSON(&Comment)
	} else {
		ctx.ShouldBind(&Comment)
	}

	commentId, _ := strconv.Atoi(ctx.Param("commentId"))

	Comment.ID = uint(commentId)

	err := db.Preload("User",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	err = db.Preload("User", func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Username", "Email", "Age")
	}).Preload("Photo",func(d *gorm.DB) *gorm.DB {
		return d.Select("ID", "Title", "Caption", "PhotoUrl")
	}).First(&Comment, commentId).Error

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Data Not Found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Comment)
}

// DeleteComment godoc
// @Summary Delete comment identified by the given Id
// @Description Delete the comment corresponding to the input Id
// @Tags comments
// @Accept json
// @Produce json
// @Param Id path int true "ID of the comment to be deleted"
// @Success 204 "No Content"
// @Router /comments/{id} [delete]
func DeleteComment(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		Comment models.Comment
	)

	commentId, _ := strconv.Atoi(ctx.Param("commentId"))

	err := db.First(&Comment, "id = ?", commentId).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found!",
		})
		return
	}
	err = db.Delete(&Comment).Error

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
			"error": "Error deleting comment data",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": fmt.Sprintf("Comment with id %v has been successfully deleted", commentId),
	})
}
