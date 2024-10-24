package models

import "time"

type Author struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:200;not null;unique"`
	Password  string `gorm:"size:200;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
