package v2

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/nuba55yo/go-101-BasicCRUD/dto"
	"github.com/nuba55yo/go-101-BasicCRUD/service"
	"github.com/gin-gonic/gin"
)

// tag แยกกับ v1 เพื่อไม่สับสนใน Swagger

// @Summary List books (v2)
// @Tags books-v2
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /books [get]
func GetBooks(svc service.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := svc.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot get books"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": "v2", "data": books})
	}
}

// @Summary Get book by id (v2)
// @Tags books-v2
// @Produce json
// @Param id path int true "book id"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func GetBook(svc service.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID, _ := strconv.Atoi(c.Param("id"))
		book, err := svc.GetByID(uint(bookID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": "v2", "data": book})
	}
}

// @Summary Create book (v2)
// @Tags books-v2
// @Accept json
// @Produce json
// @Param body body dto.CreateBookRequest true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books [post]
func CreateBook(svc service.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		created, err := svc.Create(req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrTitleExists):
				c.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(err, service.ErrBadInput):
				c.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
			}
			return
		}
		c.JSON(http.StatusCreated, gin.H{"version": "v2", "data": created})
	}
}

// @Summary Update book (v2)
// @Tags books-v2
// @Accept json
// @Produce json
// @Param id path int true "book id"
// @Param body body dto.UpdateBookRequest true "payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /books/{id} [put]
func UpdateBook(svc service.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID, _ := strconv.Atoi(c.Param("id"))
		var req dto.UpdateBookRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updated, err := svc.Update(uint(bookID), req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrTitleExists):
				c.JSON(http.StatusConflict, gin.H{"error": "title already exists"})
			case errors.Is(err, service.ErrBadInput):
				c.JSON(http.StatusBadRequest, gin.H{"error": "title and author are required"})
			default:
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"version": "v2", "data": updated})
	}
}

// @Summary Delete book (soft delete) (v2)
// @Tags books-v2
// @Param id path int true "book id"
// @Success 204
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func DeleteBook(svc service.BookService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID, _ := strconv.Atoi(c.Param("id"))
		if err := svc.Delete(uint(bookID)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
