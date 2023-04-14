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

func UserRegister(c *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(c)
		User = models.User{}
	)
	_, _ = db, contentType

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":		User.ID,
		"username": User.Username,
		"age": User.Age,
		"email": User.Email,
	})
}

func UserLogin(c *gin.Context)  {
	var (
		db = database.GetDB()
		contentType = helpers.GetContentType(c)
		User = models.User{}
	)
	_, _ = db, contentType
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password
	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "Invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}