package controllers

import (
	"fmt"
	"go-BookManagementSystem/config"
	"go-BookManagementSystem/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Exported function (uppercase 'G' to make it accessible outside the package)
func GetBooks(c *gin.Context) {
	fmt.Println("insude Get Books")
	var books []models.Book
	// Retrieve books from the database
	config.DB.Find(&books)
	// Respond with the list of books as JSON
	c.JSON(http.StatusOK, books)
}

func GetBookById(c *gin.Context) {
	id := c.Param("id")

	fmt.Printf("Value: %v, Type: %T \n", id, id)

	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not Found"}) //interface can have anythype of value
		return
	}
	c.JSON(http.StatusOK, book)

}

func PostBook(c *gin.Context) {

	idstr := c.DefaultPostForm("id", "")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID, must be integer"})
		return
	}

	name := c.DefaultPostForm("name", "")
	author := c.DefaultPostForm("author", "")
	genre := c.DefaultPostForm("genre", "")

	if name == "" || author == "" || genre == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to create new book values are empty"})
		return
	}

	book := models.Book{UID: id, NAME: name, AUTHOR: author, GENRE: genre}

	if err := config.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to create new book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book, "message": "BOOK CREATED"})

}
func PutBook(c *gin.Context) {
	fmt.Println("INSIDE PUT")
	idstr := c.DefaultPostForm("id", "")
	id, err := strconv.Atoi(idstr)
	fmt.Printf("the value is  %v and type is %T\n", id, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "id is not integer"})
		return
	}
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book Not Found"})
		return
	}
	config.DB.First(&book, id)
	book.NAME = c.DefaultPostForm("name", "")
	book.GENRE = c.DefaultPostForm("genre", "")
	book.AUTHOR = c.DefaultPostForm("author", "")
	fmt.Println(book.UID, book.NAME, book.GENRE, book.AUTHOR)

	if book.NAME == "" || book.GENRE == "" || book.AUTHOR == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to create new book values are empty"})
		return
	}

	if err := config.DB.Model(&book).Where("uid = ?", book.UID).Updates(models.Book{NAME: book.NAME, AUTHOR: book.AUTHOR, GENRE: book.GENRE}).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book was not updated"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book, "message": "BOOK Updated"})

}

func DeleteBook(c *gin.Context) {
	strid := c.Param("id")
	id, err := strconv.Atoi(strid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "id is not integer"})
		return
	}
	var book models.Book
	if err := config.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to find ID"})
		return
	}

	if err := config.DB.Delete(&book, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"book": book, "message": "BOOK Updated"})

}
