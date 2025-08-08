package dto

type CreateBookRequest struct {
	Title  string `json:"title"  binding:"required,min=1"`
	Author string `json:"author" binding:"required,min=1"`
}
type UpdateBookRequest struct {
	Title  string `json:"title"  binding:"required,min=1"`
	Author string `json:"author" binding:"required,min=1"`
}
