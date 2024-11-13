package db

import "time"

type Note struct {
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time
	NoteID    string `gorm:"primaryKey"`
	Title     string `gorm:"omitempty"`
	Body      string `gorm:"omitempty"`
}

type User struct {
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
	Username  string    `gorm:"primaryKey"`
	Email     string
	Pass      string
}
