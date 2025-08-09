package router

import (
	"github.com/gin-gonic/gin"

	v1 "github.com/nuba55yo/go-101-bookapi/http/handlers/v1"
	v2 "github.com/nuba55yo/go-101-bookapi/http/handlers/v2"
	"github.com/nuba55yo/go-101-bookapi/pkg/logger"
	"github.com/nuba55yo/go-101-bookapi/service"
)

func New(bookService service.BookService) *gin.Engine {
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Logger(), gin.Recovery(), logger.AccessLog())

	// v1 -> ต้องเรียก v1.* เท่านั้น
	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/books", v1.GetBooks(bookService))
		apiV1.GET("/books/:id", v1.GetBook(bookService))
		apiV1.POST("/books", v1.CreateBook(bookService))
		apiV1.PUT("/books/:id", v1.UpdateBook(bookService))
		apiV1.DELETE("/books/:id", v1.DeleteBook(bookService))
	}

	// v2 -> ต้องเรียก v2.* เท่านั้น
	apiV2 := r.Group("/api/v2")
	{
		apiV2.GET("/books", v2.GetBooks(bookService))
		apiV2.GET("/books/:id", v2.GetBook(bookService))
		apiV2.POST("/books", v2.CreateBook(bookService))
		apiV2.PUT("/books/:id", v2.UpdateBook(bookService))
		apiV2.DELETE("/books/:id", v2.DeleteBook(bookService))
	}
	return r
}
