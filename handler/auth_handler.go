package handler

import (
	"go-todo-api/config"
	"go-todo-api/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secret_key")

func Register(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 🔥 HASH PASSWORD
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)

	config.DB.Create(&user)

	c.JSON(201, gin.H{
		"message": "User registered",
	})
}

func Login(c *gin.Context) {
	var input model.User
	var user model.User

	c.ShouldBindJSON(&input)

	config.DB.Where("username = ?", input.Username).First(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// 🔹 access token (short)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	accessString, _ := accessToken.SignedString(jwtKey)

	// 🔹 refresh token (long)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshString, _ := refreshToken.SignedString(jwtKey)

	// simpan ke DB
	user.RefreshToken = refreshString
	config.DB.Save(&user)

	c.JSON(200, gin.H{
		"access_token":  accessString,
		"refresh_token": refreshString,
	})
}

func RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user model.User
	config.DB.Where("refresh_token = ?", input.RefreshToken).First(&user)

	if user.ID == 0 {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// validasi token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	// generate access token baru
	newAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	newAccessString, _ := newAccess.SignedString(jwtKey)

	c.JSON(200, gin.H{
		"access_token": newAccessString,
	})
}
