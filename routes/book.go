package routes

import (
	"BookManagemantGoLang/database"
	"BookManagemantGoLang/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SaveBook(c *gin.Context) {
	var book models.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	authorID, exists := c.Get("authorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	book.AuthorID = authorID.(uint)

	if err := database.DB.Create(&book).Error; err != nil {
		log.Println("Could not create book:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create book"})
		return
	}

	// Preload the associated author (optional, to return the full book object with author)
	if err := database.DB.Preload("Author").First(&book).Error; err != nil {
		log.Println("Could not find the book:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book saved successfully", "book": book})
}

func GetAllBooks(c *gin.Context) {
	var books []models.Book

	if err := database.DB.Preload("Author").Find(&books).Error; err != nil {
		log.Println("Could not retrieve books:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve books"})
		return
	}

	if len(books) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No books found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get books successfully", "books": books})
}

func GetBook(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	authorID, exists := c.Get("authorID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	book.AuthorID = authorID.(uint) // Set the authorID from the token claims

	if err := database.DB.Preload("Author").First(&book, id).Error; err != nil {
		log.Println("Could not find the book:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get book successfully", "book": book})
}

// func GetBookByAuthorID(c *gin.Context) {
// 	var author models.Author
// 	id := c.Param("id") // get by id

// 	if err := database.DB.Find(&author, id).Error; err != nil {
// 		log.Println("Could not find the book:", err)
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Get book by authorID successfully", "author": author})
// }

func GetBookByAuthorID(c *gin.Context) {
	var books []models.Book
	authorID := c.Param("id")

	if err := database.DB.Where("author_id = ?", authorID).Find(&books).Error; err != nil {
		log.Println("Could not find books for the author:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Books not found for this author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Books retrieved successfully", "books": books})
}
