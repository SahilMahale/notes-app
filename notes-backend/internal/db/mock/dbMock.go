package mock

import (
	"os"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestDB(t *testing.T) db.DbConnection {
	tmpfile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	})

	gdb, err := gorm.Open(sqlite.Open(tmpfile.Name()), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	err = gdb.AutoMigrate(db.Note{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}
	err = gdb.AutoMigrate(db.User{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db.DbConnection{Db: gdb}
}
