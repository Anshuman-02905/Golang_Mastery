package routes

import (
	"bookstore/controllers"

	"github.com/gin-gonic/gin"
)

// CREATE ROUTES
//  GetAllBooks()
// CreateBook()
// GetBookByID()
// UpdateBook()
// DeleteBook()

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.GET("/book/", controllers.GetAllBooksFunc)
	r.POST("/book/", controllers.CreateBookFunc)
	r.GET("/book/:id", controllers.GetBookByIDFunc)
	r.PUT("/book/:id", controllers.UpdateBookFunc)
	r.DELETE("/book/:id", controllers.UpdateBookFunc)

	return r

}
