package models

type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`	
	Price       string  `json:"price"`
  Cover       string  `json:"cover"`
}
