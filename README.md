# CRUD Book API ด้วย Go + Gin + GORM + PostgreSQL

มี **Swagger**, **API Version Control** , **Validation (ห้ามชื่อซ้ำ)** และ **Logging** (บันทึก request/response, แยกโฟลเดอร์รายวัน, หมุนไฟล์ใหม่ทุก 10 นาที)

## Go packages

```bash
# เว็บเฟรมเวิร์ก + ORM + Postgres driver + .env
go get github.com/gin-gonic/gin gorm.io/gorm gorm.io/driver/postgres github.com/joho/godotenv

# Swagger UI (router ใช้)
go get github.com/swaggo/gin-swagger github.com/swaggo/files

# เครื่องมือ gen เอกสาร (รันครั้งเดียวพอ)
go install github.com/swaggo/swag/cmd/swag@latest

# เก็บ dependency ให้เรียบร้อย
go mod tidy

# สร้างเอกสาร Swagger (ทุกครั้งที่เพิ่ม/แก้ annotations ใน handlers)
swag init -g main.go
```

## โครงสร้างโปรเจกต์

```text
database/           # เชื่อมต่อ DB (GORM)
docs/               # Swagger (gen โดย swag)
dto/                # Request DTO
http/
  handlers/         # Controller: รับ/ส่ง HTTP (ไม่มี business logic)
  router/           # gin engine + middleware + routes
models/             # GORM models
pkg/logger/         # Middleware + file-rotate logs (ทุก 10 นาที)
repository/         # Data access (GORM)
service/            # Business logic/validation
main.go             # จุดเริ่มโปรแกรม (DI + Run)
```

## SQL: สร้างตาราง `books` (PostgreSQL)

```sql
CREATE TABLE IF NOT EXISTS public.books (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT        NOT NULL,
    author      TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ NULL
);

