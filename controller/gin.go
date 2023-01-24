package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/olusamimaths/go-jwt/model"
)

func Signup(c *gin.Context) {
	var reqUser model.User

	// Bind incoming user Data
	if err := c.ShouldBind(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Check existence of the user in DB
	var dbUser model.User
	model.DB.Where("email = ?", reqUser.Email).First(&dbUser)
	if dbUser.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user with email found, please login",
		})
		return
	}

	err := reqUser.GeneratePasswordHash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error occured",
		})
		return
	}

	res := model.DB.Create(&reqUser)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"user": reqUser,
	})
}

func Login(c *gin.Context) {
	var reqUser model.User
	// Bind incoming user Data
	if err := c.ShouldBind(&reqUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// Check existence of the user in DB
	var dbUser model.User
	model.DB.Where("email = ?", reqUser.Email).First(&dbUser)
	if dbUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if dbUser.CheckPasswordHash(reqUser.Password) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": dbUser.Email,
			"exp": time.Now().Add(time.Minute * 10).Unix(),
		})

		// sign and get the encode d token as a string
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"user": dbUser,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
	}
}

func Resources(c *gin.Context) {
	var users []model.User
	res := model.DB.Find(&users)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error fetching users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
