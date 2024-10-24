package main

import (
	"BookManagemantGoLang/auth"
	"BookManagemantGoLang/database"
	"BookManagemantGoLang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()

	r := gin.Default()

	r.POST("/register", routes.Register)
	r.POST("/login", routes.Login)
	r.GET("/getAuthors", routes.GetAllAuthors)

	r.POST("/saveNewBooks", auth.AuthMiddleware(), routes.SaveBook)
	r.GET("/getAllBooks/", auth.AuthMiddleware(), routes.GetAllBooks)
	r.GET("/getBooksByID/:id", auth.AuthMiddleware(), routes.GetBook)
	r.GET("/getBookByAuthorID/:id", auth.AuthMiddleware(), routes.GetBookByAuthorID)

	r.POST("/addInventory", auth.AuthMiddleware(), routes.AddInventoryForBook)
	r.POST("/getAllInventoryOfEachAuthor", auth.AuthMiddleware(), routes.GetAllInventoryOfEachAuthor)

	r.POST("/orderBook", auth.AuthMiddleware(), routes.OrderBook)
	r.Run(":8080")
}
