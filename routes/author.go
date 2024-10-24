package routes

import (
	"BookManagemantGoLang/auth"
	"BookManagemantGoLang/database"
	"BookManagemantGoLang/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(author.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Could not hash password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	author.Password = string(hashedPassword)

	if err := database.DB.Create(&author).Error; err != nil {
		log.Println("Could not create author:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create author"})
		return
	}

	token, err := auth.GenerateJWT(&author)
	if err != nil {
		log.Println("Could not generate token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.String(http.StatusOK, token)
	// c.JSON(http.StatusOK, gin.H{token})
}

func Login(c *gin.Context) {
	var author models.Author
	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var dbAuthor models.Author
	if err := database.DB.Where("username = ?", author.Username).First(&dbAuthor).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbAuthor.Password), []byte(author.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateJWT(&dbAuthor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.String(http.StatusOK, token)
	// c.JSON(http.StatusOK, gin.H{token})
}

// GetAllAuthors retrieves all authors from the database
func GetAllAuthors(c *gin.Context) {
	var authors []models.Author
	if err := database.DB.Find(&authors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get authors"})
		return
	}
	c.JSON(http.StatusOK, authors)
}
