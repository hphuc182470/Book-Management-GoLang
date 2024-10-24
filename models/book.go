package models

import "time"

type Book struct {
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"size:200;not null"`
	AuthorID      uint   `gorm:"not null"`
	Author        Author `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"`
	PublishedYear int
	ISBN          string `gorm:"size:20;unique"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
