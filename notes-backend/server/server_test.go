package server

import (
	"reflect"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db"
)

func TestNewNotesService(t *testing.T) {
	type args struct {
		appname string
		ip      string
		db      db.DbConnection
	}
	tests := []struct {
		name string
		args args
		want notesService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotesService(tt.args.appname, tt.args.ip, tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotesService() = %v, want %v", got, tt.want)
			}
		})
	}
}
