package repository

import (
	"strings"
	"time"

	"github.com/nuba55yo/go-101-BasicCRUD/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BookRepository สัญญาให้ service เรียกใช้งาน
type BookRepository interface {
	Create(book *models.Book) error
	GetAll() ([]models.Book, error)
	GetByID(bookID uint) (*models.Book, error)
	Update(book *models.Book) error
	SoftDelete(bookID uint) error
	ExistsActiveByTitle(title string) (bool, error)
	ExistsActiveByTitleExceptID(title string, bookID uint) (bool, error)
}

type bookRepository struct{ db *gorm.DB }

// NewBookRepository รับ *gorm.DB และคืน Repository ที่พร้อมใช้งาน
func NewBookRepository(database *gorm.DB) BookRepository { return &bookRepository{db: database} }

func (repository *bookRepository) Create(book *models.Book) error {
	return repository.db.Create(book).Error
}

func (repository *bookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book
	err := repository.db.Where("deleted_at IS NULL").Order("id DESC").Find(&books).Error
	return books, err
}

func (repository *bookRepository) GetByID(bookID uint) (*models.Book, error) {
	var book models.Book
	err := repository.db.Where("id = ? AND deleted_at IS NULL", bookID).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (repository *bookRepository) Update(book *models.Book) error {
	book.UpdatedAt = time.Now()
	return repository.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Clauses(clause.Returning{}).
		Save(book).Error
}

func (repository *bookRepository) SoftDelete(bookID uint) error {
	return repository.db.Model(&models.Book{}).
		Where("id = ?", bookID).
		Update("deleted_at", gorm.Expr("now()")).Error
}

// ตรวจชื่อซ้ำ (ไม่แคร์ตัวพิมพ์)
func (repository *bookRepository) ExistsActiveByTitle(title string) (bool, error) {
	normalized := strings.TrimSpace(title)
	var count int64
	err := repository.db.Model(&models.Book{}).
		Where("deleted_at IS NULL AND lower(title)=lower(?)", normalized).
		Count(&count).Error
	return count > 0, err
}

func (repository *bookRepository) ExistsActiveByTitleExceptID(title string, bookID uint) (bool, error) {
	normalized := strings.TrimSpace(title)
	var count int64
	err := repository.db.Model(&models.Book{}).
		Where("deleted_at IS NULL AND id <> ? AND lower(title)=lower(?)", bookID, normalized).
		Count(&count).Error
	return count > 0, err
}
