package routes

import (
	"BookManagemantGoLang/database"
	"BookManagemantGoLang/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderBook(c *gin.Context) {
	// begin transactions
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not start transaction"})
		return
	}

	// bind json
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		tx.Rollback()
		return
	}

	var inventory models.Inventory
	if err := tx.Where("book_id = ?", order.BookID).First(&inventory).Error; err != nil {
		log.Println("Inventory not found for the book:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inventory not found"})
		tx.Rollback()
		return
	}

	// Check quantity enough
	if inventory.Quantity < order.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough inventory available"})
		tx.Rollback()
		return
	}

	// Create the order
	if err := tx.Create(&order).Error; err != nil {
		log.Println("Could not create order:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create order"})
		tx.Rollback()
		return
	}

	// minus quantity
	inventory.Quantity -= order.Quantity
	if err := tx.Save(&inventory).Error; err != nil {
		log.Println("Could not update inventory:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update inventory"})
		tx.Rollback()
		return
	}

	var createdOrder models.Order
	if err := tx.Preload("Book.Author").First(&createdOrder, order.ID).Error; err != nil {
		log.Println("Could not preload order with book and author:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not preload order with book and author"})
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Transaction commit failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "order": createdOrder})
}
