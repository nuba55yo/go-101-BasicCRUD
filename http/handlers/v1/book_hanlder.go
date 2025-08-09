package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nuba55yo/go-101-BasicCRUD/dto"
	"github.com/nuba55yo/go-101-BasicCRUD/service"

	"github.com/gin-gonic/gin"
)

// @Summary ดึงรายการหนังสือทั้งหมด
// @Tags books
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /books [get]
func GetBooks(bookService service.BookService) gin.HandlerFunc {
	return func(context *gin.Context) {
		books, err := bookService.GetAll()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get books"})
			return
		}
		context.JSON(http.StatusOK, books)
	}
}

// @Summary ดึงหนังสือตามรหัส
// @Tags books
// @Produce json
// @Param id path int true "book id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func GetBook(bookService service.BookService) gin.HandlerFunc {
	return func(context *gin.Context) {
		bookID, _ := strconv.Atoi(context.Param("id"))
		book, err := bookService.GetByID(uint(bookID))
		if err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		context.JSON(http.StatusOK, book)
	}
}

// @Summary สร้างหนังสือใหม่
// @Tags books
// @Accept json
// @Produce json
// @Param body body dto.CreateBookRequest true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books [post]
func CreateBook(bookService service.BookService) gin.HandlerFunc {
	return func(context *gin.Context) {
		var requestBody dto.CreateBookRequest
		if err := context.ShouldBindJSON(&requestBody); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		createdBook, err := bookService.Create(requestBody)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrTitleExists):
				// 409 สำหรับเคสชื่อซ้ำ
				context.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(err, service.ErrBadInput):
				context.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			default:
				context.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			}
			return
		}
		context.JSON(http.StatusCreated, createdBook)
	}
}

// @Summary แก้ไขหนังสือ
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "book id"
// @Param body body dto.UpdateBookRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books/{id} [put]
func UpdateBook(bookService service.BookService) gin.HandlerFunc {
	return func(context *gin.Context) {
		bookID, _ := strconv.Atoi(context.Param("id"))

		var requestBody dto.UpdateBookRequest
		if err := context.ShouldBindJSON(&requestBody); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updatedBook, err := bookService.Update(uint(bookID), requestBody)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrTitleExists):
				context.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(err, service.ErrBadInput):
				context.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			default:
				context.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			}
			return
		}
		context.JSON(http.StatusOK, updatedBook)
	}
}

// @Summary ลบหนังสือ (soft delete)
// @Tags books
// @Param id path int true "book id"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func DeleteBook(bookService service.BookService) gin.HandlerFunc {
	return func(context *gin.Context) {
		bookID, _ := strconv.Atoi(context.Param("id"))
		if err := bookService.Delete(uint(bookID)); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		context.Status(http.StatusNoContent)
	}
}
