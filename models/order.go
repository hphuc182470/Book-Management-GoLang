package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	BookID    uint      `gorm:"not null"`
	Book      Book      `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	OrderDate time.Time `gorm:"default:current_timestamp"`
	Quantity  int       `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
