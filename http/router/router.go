package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/nuba55yo/go-101-bookapi/http/handlers"
	"github.com/nuba55yo/go-101-bookapi/pkg/logger"
	"github.com/nuba55yo/go-101-bookapi/service"
)

// New สร้าง *gin.Engine พร้อม middleware และทุกเส้นทางของระบบ
func New(bookService service.BookService) *gin.Engine {
	// ใช้ gin.New() เพื่อกำหนด middleware เอง
	router := gin.New()

	// กันเตือน proxy ใน dev และเพิ่ม middleware พื้นฐาน
	_ = router.SetTrustedProxies(nil)
	router.Use(gin.Logger())    // access log แบบสั้นๆ ไป stdout
	router.Use(gin.Recovery())  // กันโปรแกรมตายเมื่อ panic

	// เขียน log แบบละเอียดของเรา (เก็บทั้ง request + response, ทุกสถานะ)
	router.Use(logger.AccessLog())

	// Swagger UI (อย่าไปประกาศซ้ำที่ main)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// เส้นทางของ Book
	router.GET("/books", handlers.GetBooks(bookService))
	router.GET("/books/:id", handlers.GetBook(bookService))
	router.POST("/books", handlers.CreateBook(bookService))
	router.PUT("/books/:id", handlers.UpdateBook(bookService))
	router.DELETE("/books/:id", handlers.DeleteBook(bookService))

	return router
}
