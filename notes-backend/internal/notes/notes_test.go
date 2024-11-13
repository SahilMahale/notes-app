package notes

import (
	"reflect"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db/mock"
	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/models"
)

func TestNotesController_CreateNote(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	nCtrl := NewNoteController(dbMock)
	type args struct {
		title string
		body  string
	}
	tests := []struct {
		name      string
		args      args
		wantMyerr helper.MyHTTPErrors
	}{
		{
			name: "CreateNote: Positive case",
			args: args{
				title: "Test Note",
				body:  "Test body",
			},
			wantMyerr: helper.MyHTTPErrors{
				Err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := nCtrl.CreateNote(tt.args.title, tt.args.body)
			if got.Err != tt.wantMyerr.Err {
				t.Errorf("NotesController.CreateNote() gotMyERR = %v, wantMyErr %v", got.Err, tt.wantMyerr.Err)
			}
		})
	}
}

func TestNotesController_UpdateNote(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	nCtrl := NewNoteController(dbMock)
	noteID, err := nCtrl.CreateNote("Test Note", "Test Body")
	if err.Err != nil {
		t.Error(err.Err)
	}
	type args struct {
		noteID string
		note   models.NotePatchRequest
	}
	tests := []struct {
		name      string
		args      args
		want      models.NotePatchResp
		wantMyerr helper.MyHTTPErrors
	}{
		{
			name: "UpdateNote: Positive Case",
			args: args{
				noteID: noteID,
				note: models.NotePatchRequest{
					Body: "Test Body updated",
				},
			},
			want: models.NotePatchResp{
				NoteID: noteID,
				Title:  "Test Note",
				Body:   "Test Body updated",
			},
			wantMyerr: helper.MyHTTPErrors{
				Err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotMyErr := nCtrl.UpdateNote(tt.args.noteID, tt.args.note)
			if got.Body != tt.want.Body && got.NoteID != tt.want.NoteID {
				t.Errorf("NotesController.UpdateNote() got = %v, want %v", got, tt.want)
			}
			if gotMyErr.Err != tt.wantMyerr.Err {
				t.Errorf("NotesController.UpdateNote() got1 = %v, want %v", gotMyErr, tt.wantMyerr)
			}
		})
	}
}

func TestNotesController_GetAllNotes(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	nCtrl := NewNoteController(dbMock)
	noteID1, err := nCtrl.CreateNote("Test Note", "Test Body")
	if err.Err != nil {
		t.Error(err.Err)
	}
	noteID2, err := nCtrl.CreateNote("Test Note2", "Test Body2")
	if err.Err != nil {
		t.Error(err.Err)
	}
	tests := []struct {
		want    models.NotesResp
		wantErr helper.MyHTTPErrors
		name    string
	}{
		{
			name: "GetAllNotes: Positive",
			want: models.NotesResp{
				{
					NoteID: noteID1,
					Title:  "Test Note",
					Body:   "Test Body",
				},
				{
					NoteID: noteID2,
					Title:  "Test Note2",
					Body:   "Test Body2",
				},
			},
			wantErr: helper.MyHTTPErrors{
				Err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := nCtrl.GetAllNotes()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotesController.GetAllNotes() got = %v, want %v", got, tt.want)
			}
			if gotErr.Err != tt.wantErr.Err {
				t.Errorf("NotesController.GetAllNotes() got1 = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
