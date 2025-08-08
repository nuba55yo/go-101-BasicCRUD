package main

import (
	"log"
	"os"

	// Swagger docs (จะมีโฟลเดอร์ docs หลังรัน `swag init`)
	_ "github.com/nuba55yo/go-101-bookapi/docs"

	"github.com/nuba55yo/go-101-bookapi/database"
	"github.com/nuba55yo/go-101-bookapi/http/router"
	"github.com/nuba55yo/go-101-bookapi/models"
	"github.com/nuba55yo/go-101-bookapi/pkg/logger"
	"github.com/nuba55yo/go-101-bookapi/repository"
	"github.com/nuba55yo/go-101-bookapi/service"
)

// @title Book API (Go)
// @version 1.0
// @description ตัวอย่าง API จัดการหนังสือ (Gin + GORM + Postgres)
// @schemes http
// @host localhost:8080
// @BasePath /
func main() {
	// ปิดไฟล์ log ให้เรียบร้อยตอนโปรเซสจบ
	defer logger.Close()

	// 1) เชื่อมต่อฐานข้อมูล + AutoMigrate สร้าง/อัปเดตตารางตาม model
	database.Connect()
	if err := database.DB.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal(err)
	}

	// 2) ประกอบ dependencies แบบ manual (Repository → Service → Router)
	bookRepository := repository.NewBookRepository(database.DB) // คุยกับ DB โดยตรง
	bookService := service.NewBookService(bookRepository)       // ธุรกิจ/กฎเกณฑ์
	httpRouter := router.New(bookService)                        // เส้นทาง HTTP ทั้งหมด

	// 3) อ่านพอร์ตจาก ENV (มีค่า default = 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 4) เปิดเว็บเซิร์ฟเวอร์
	_ = httpRouter.Run(":" + port)
}
