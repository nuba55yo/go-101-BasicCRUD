package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// สำคัญ: ต้อง gen เอกสารไว้ก่อน และ import สองเวอร์ชันนี้
	_ "github.com/nuba55yo/go-101-bookapi/docs/v1"
	_ "github.com/nuba55yo/go-101-bookapi/docs/v2"

	"github.com/nuba55yo/go-101-bookapi/database"
	"github.com/nuba55yo/go-101-bookapi/http/router"
	"github.com/nuba55yo/go-101-bookapi/models"
	"github.com/nuba55yo/go-101-bookapi/pkg/logger"
	"github.com/nuba55yo/go-101-bookapi/repository"
	"github.com/nuba55yo/go-101-bookapi/service"
)

func main() {
	defer logger.Close()

	// DB + AutoMigrate
	database.Connect()
	if err := database.DB.AutoMigrate(&models.Book{}); err != nil {
		log.Fatal(err)
	}

	// DI
	bookRepo := repository.NewBookRepository(database.DB)
	bookSvc := service.NewBookService(bookRepo)
	r := router.New(bookSvc)

	// ---------- เสิร์ฟสเปค (doc.json) แยกเวอร์ชัน ----------
	// อย่าลบ InstanceName ออก เพื่อแยก v1/v2 ให้ชัดเจน
	r.GET("/docs/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("v1")))
	r.GET("/docs/v2/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName("v2")))

	// ---------- Swagger UI แบบหน้าเดียว + Dropdown v1/v2 ----------
	// ไม่ใช้ ginSwagger.URL/URLs เลย เพื่อเลี่ยงปัญหาเวอร์ชัน/แคช
	r.GET("/swagger", swaggerIndex())
	r.GET("/swagger/index.html", func(c *gin.Context) {
		// ให้พารามิเตอร์ที่มากับ /swagger/index.html ติดไปด้วย
		raw := c.Request.URL.RawQuery
		if raw != "" {
			raw = "?" + raw
		}
		c.Redirect(302, "/swagger"+raw)
	})

	// Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	_ = r.Run(":" + port)
}

// swaggerIndex คืน HTML ของ Swagger UI (ใช้ CDN) และมี dropdown v1/v2
func swaggerIndex() gin.HandlerFunc {
	const html = `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <meta name="viewport" content="width=device-width, initial-scale=1"/>
  <style> body{margin:0; padding:0;} </style>
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <!-- ต้องมีไฟล์นี้เพิ่ม เพื่อใช้ StandaloneLayout -->
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    const params  = new URLSearchParams(location.search);
    const primary = params.get('urls.primaryName') || 'v1';

    window.ui = SwaggerUIBundle({
      dom_id: '#swagger-ui',
      urls: [
        { url: '/docs/v1/doc.json', name: 'v1' },
        { url: '/docs/v2/doc.json', name: 'v2' }
      ],
      'urls.primaryName': primary,
      deepLinking: true,
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: 'StandaloneLayout'
    });
  </script>
</body>
</html>`
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, html)
	}
}
