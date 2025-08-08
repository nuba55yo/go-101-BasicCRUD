package models

import "time"

type Book struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Title     string     `json:"title" gorm:"not null"`
	Author    string     `json:"author" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
