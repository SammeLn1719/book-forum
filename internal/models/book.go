package models

type Book struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Author      string  `json:"authar"`
	Description string  `json:"description"`
	Price       float64 `json:"prise"`
}
