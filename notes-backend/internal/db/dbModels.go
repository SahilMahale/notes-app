package db

import "time"

type Note struct {
	NoteID    string    `gorm:"primaryKey"`
	Title     string    `gorm:"omitempty"`
	Body      string    `gorm:"omitempty"`
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time
}

type User struct {
	Username  string `gorm:"primaryKey"`
	Email     string
	Pass      string
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}
