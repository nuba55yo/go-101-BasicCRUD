package service

import (
	"errors"
	"strings"

	"github.com/nuba55yo/go-101-BasicCRUD/dto"
	"github.com/nuba55yo/go-101-BasicCRUD/models"
	"github.com/nuba55yo/go-101-BasicCRUD/pkg/logger"
	"github.com/nuba55yo/go-101-BasicCRUD/repository"
	"gorm.io/gorm"
)

// error ธุรกิจที่ handler จะใช้ตัดสินใจแปลงเป็นสถานะ HTTP
var (
	ErrTitleExists = errors.New("title already exists")
	ErrBadInput    = errors.New("invalid input")
)

// BookService กำหนดสัญญาให้เลเยอร์บนเรียกใช้งาน
type BookService interface {
	Create(dto.CreateBookRequest) (*models.Book, error)
	GetAll() ([]models.Book, error)
	GetByID(bookID uint) (*models.Book, error)
	Update(bookID uint, request dto.UpdateBookRequest) (*models.Book, error)
	Delete(bookID uint) error
}

// bookService โครงสร้างภายใน (ซ่อนหลัง interface) — หลีกเลี่ยงใช้ตัวอักษรเดียว
type bookService struct {
	repository repository.BookRepository
}

// NewBookService คืน service พร้อม repository ที่ถูกฉีดเข้ามา
func NewBookService(bookRepository repository.BookRepository) BookService {
	return &bookService{repository: bookRepository}
}

// normalize ตัดช่องว่างหัว-ท้าย เพื่อกันเคสส่ง "  ชื่อ  "
func normalize(title, author string) (string, string) {
	return strings.TrimSpace(title), strings.TrimSpace(author)
}

func (serviceImpl *bookService) Create(request dto.CreateBookRequest) (*models.Book, error) {
	title, author := normalize(request.Title, request.Author)
	if title == "" || author == "" {
		return nil, ErrBadInput
	}

	// ตรวจชื่อซ้ำ (ไม่สนตัวพิมพ์เล็ก/ใหญ่)
	exists, err := serviceImpl.repository.ExistsActiveByTitle(title)
	if err != nil {
		logger.Errorf("books", "check duplicate failed: %v", err)
		return nil, err
	}
	if exists {
		return nil, ErrTitleExists
	}

	newBook := &models.Book{Title: title, Author: author}
	if err := serviceImpl.repository.Create(newBook); err != nil {
		logger.Errorf("books", "create failed: %v", err)
		return nil, err
	}

	logger.Infof("books", "created id=%d title=%s", newBook.ID, newBook.Title)
	return newBook, nil
}

func (serviceImpl *bookService) GetAll() ([]models.Book, error) {
	books, err := serviceImpl.repository.GetAll()
	if err != nil {
		logger.Errorf("books", "list failed: %v", err)
	}
	return books, err
}

func (serviceImpl *bookService) GetByID(bookID uint) (*models.Book, error) {
	book, err := serviceImpl.repository.GetByID(bookID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		logger.Errorf("books", "get failed: %v", err)
	}
	return book, err
}

func (serviceImpl *bookService) Update(bookID uint, request dto.UpdateBookRequest) (*models.Book, error) {
	title, author := normalize(request.Title, request.Author)
	if title == "" || author == "" {
		return nil, ErrBadInput
	}

	book, err := serviceImpl.repository.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	// ตรวจชื่อซ้ำ ยกเว้นเล่มตัวเอง
	exists, err := serviceImpl.repository.ExistsActiveByTitleExceptID(title, bookID)
	if err != nil {
		logger.Errorf("books", "check duplicate failed: %v", err)
		return nil, err
	}
	if exists {
		return nil, ErrTitleExists
	}

	book.Title = title
	book.Author = author

	if err := serviceImpl.repository.Update(book); err != nil {
		logger.Errorf("books", "update failed: %v", err)
		return nil, err
	}

	logger.Infof("books", "updated id=%d title=%s", book.ID, book.Title)
	return book, nil
}

func (serviceImpl *bookService) Delete(bookID uint) error {
	if err := serviceImpl.repository.SoftDelete(bookID); err != nil {
		logger.Errorf("books", "delete failed id=%d: %v", bookID, err)
		return err
	}
	logger.Infof("books", "deleted id=%d", bookID)
	return nil
}
