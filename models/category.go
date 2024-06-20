package models

type Category struct {
	Id    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Posts []Post `json:"posts"`
}
