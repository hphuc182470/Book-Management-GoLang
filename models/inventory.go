package models

import "time"

type Inventory struct {
	ID        uint `gorm:"primaryKey"`
	BookID    uint `gorm:"not null"`
	Book      Book `gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Quantity  int  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
