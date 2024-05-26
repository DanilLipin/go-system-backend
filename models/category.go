package models

type Category struct {
	Id   uint   `json:"id" db:"id"`
	Name string `json:"title" db:"title"`
}
