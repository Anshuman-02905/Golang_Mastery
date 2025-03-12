package models

import "time"

type Book struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"username" gorm:"column:username"`
	AuthorName string    `json:"author" gorm:"column:author"`
	Genre      string    `json:"genre" gorm:"column:genre"`
	Rent       string    `json:"rent" gorm:"column:rent"`
	Currency   string    `json:"currency" gorm:"column:currency"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"` // Auto-set on insert
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Auto-set on update
}
