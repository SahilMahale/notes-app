package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbConnection struct {
	Db *gorm.DB
}

func NewDbConnection() (DbConnection, error) {
	db, err := gorm.Open(sqlite.Open("Notes.DB"), &gorm.Config{})
	if err != nil {
		return DbConnection{}, fmt.Errorf("Error establishing DB connection: %V", err)
	}
	db.AutoMigrate(Note{})
	return DbConnection{Db: db}, nil
}
