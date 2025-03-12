package routes

import (
	"go-BookManagementSystem/controllers" // Import the controllers package

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the routes for the application
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// Correctly call the exported GetBooks function
	r.GET("/books", controllers.GetBooks)
	r.GET("/books/:id", controllers.GetBookById)
	r.POST("/books", controllers.PostBook)
	r.PUT("/books/:id", controllers.PutBook)
	r.DELETE("books/:id", controllers.DeleteBook)
	return r
}
