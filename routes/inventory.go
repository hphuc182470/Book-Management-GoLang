package routes

import (
	"BookManagemantGoLang/database"
	"BookManagemantGoLang/models"
	"log"
	"net/http"

	// "time"

	"github.com/gin-gonic/gin"
)

func AddInventoryForBook(c *gin.Context) {
	var inventory models.Inventory

	// bind json
	if err := c.ShouldBindJSON(&inventory); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// check exist
	var book models.Book
	if err := database.DB.First(&book, inventory.BookID).Error; err != nil {
		log.Println("Book not found:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	// Get authorID from the request
	authorID, exists := c.Get("authorID")
	if !exists {
		log.Println("Unauthorized: Author ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Validate
	if authorID.(uint) != book.AuthorID {
		log.Println("Author ID does not match the book's author")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Author ID does not match the book's author"})
		return
	}

	// create
	var existingInventory models.Inventory
	if err := database.DB.Where("book_id = ?", inventory.BookID).First(&existingInventory).Error; err == nil {
		// exists -> update
		existingInventory.Quantity += inventory.Quantity
		if err := database.DB.Save(&existingInventory).Error; err != nil {
			log.Println("Could not update inventory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update inventory"})
			return
		}
		// load book and author
		if err := database.DB.Preload("Book.Author").First(&existingInventory, existingInventory.ID).Error; err != nil {
			log.Println("Could not preload inventory with book and author:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not preload inventory with book and author"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Inventory updated successfully", "inventory": existingInventory})
	} else {
		// doesn't exist -> create
		if err := database.DB.Create(&inventory).Error; err != nil {
			log.Println("Could not create inventory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create inventory"})
			return
		}
		// load book and author
		if err := database.DB.Preload("Book.Author").First(&inventory, inventory.ID).Error; err != nil {
			log.Println("Could not preload inventory with book and author:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not preload inventory with book and author"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Inventory added successfully", "inventory": inventory})
	}
}

func GetAllInventoryOfEachAuthor(c *gin.Context) {
	var inventory []models.Inventory

	//if err := database.DB.Find(&inventory).Error; err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get authors"})
	//	return
	//}

	if err := database.DB.Preload("Book.Author").Find(&inventory).Error; err != nil {
		log.Println("Could not preload order with book and author:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not preload order with book and author"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}
