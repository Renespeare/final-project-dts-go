package controllers

import (
	"net/http"
	"final-project-dts-go/database"
	"final-project-dts-go/helpers"
	"final-project-dts-go/models"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

// UserRegister godoc
// @Summary Register new user
// @Description Regiter new user
// @Tags users
// @Accept json
// @Produce json
// @Param models.User body models.User true "user register"
// @Success 201 {object} models.User
// @Router /users/register [post]
func UserRegister(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(ctx)
		User = models.User{}
	)
	_, _ = db, contentType

	if contentType == appJSON {
		ctx.ShouldBindJSON(&User)
	} else {
		ctx.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":		User.ID,
		"username": User.Username,
		"age": User.Age,
		"email": User.Email,
	})
}

// UserLogin godoc
// @Summary Login user
// @Description User can login with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param models.User body models.User true "user login"
// @Success 200 {object} models.User
// @Router /users/login [post]
func UserLogin(ctx *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(ctx)
		User = models.User{}
	)
	_, _ = db, contentType
	password := ""

	if contentType == appJSON {
		ctx.ShouldBindJSON(&User)
	} else {
		ctx.ShouldBind(&User)
	}

	password = User.Password
	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}