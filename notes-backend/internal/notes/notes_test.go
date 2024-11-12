package notes

import (
	"reflect"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/db/mock"
	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/models"
)

func TestNotesController_UpdateNote(t *testing.T) {
	type fields struct {
		DbInterface db.DbConnection
	}
	type args struct {
		noteID string
		note   models.NotePatchRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   models.NotePatchResp
		want1  helper.MyHTTPErrors
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NotesController{
				DbInterface: tt.fields.DbInterface,
			}
			got, got1 := n.UpdateNote(tt.args.noteID, tt.args.note)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotesController.UpdateNote() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NotesController.UpdateNote() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestNotesController_CreateNote(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	nCtrl := NewNoteController(dbMock)
	type args struct {
		title string
		body  string
	}
	tests := []struct {
		name        string
		args        args
		wantdNoteId string
		wantMyerr   helper.MyHTTPErrors
	}{
		{
			name: "CreateNote: Positive case",
			args: args{
				title: "Test Note",
				body:  "Test body",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := nCtrl.CreateNote(tt.args.title, tt.args.body)
			if got != tt.wantdNoteId {
				t.Errorf("NotesController.CreateNote() got = %v, want %v", got, tt.wantdNoteId)
			}
			if !reflect.DeepEqual(got1, tt.wantMyerr) {
				t.Errorf("NotesController.CreateNote() got1 = %v, want %v", got1, tt.wantMyerr)
			}
		})
	}
}
