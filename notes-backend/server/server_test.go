package server

import (
	"reflect"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/gofiber/fiber/v2"
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

func Test_notesService_GetNotes(t *testing.T) {
	type fields struct {
		app         *fiber.App
		DbInterface db.DbConnection
		ip          string
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			B := &notesService{
				app:         tt.fields.app,
				DbInterface: tt.fields.DbInterface,
				ip:          tt.fields.ip,
			}
			if err := B.GetNotes(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("notesService.GetNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
