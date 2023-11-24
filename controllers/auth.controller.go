package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cocoza4/data_microservices/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
}

func NewAuthController() AuthController {
	return AuthController{}
}

func login(ctx *gin.Context) {
	var user models.User
	if ctx.BindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to read body"})
		return
	}

	// normally we'd lookup the database to get credentials and compute but since this is out of scope of this test, i will just use mock user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(), // session expires in 1 hour
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("err", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to create token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ctr *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", login)
}
