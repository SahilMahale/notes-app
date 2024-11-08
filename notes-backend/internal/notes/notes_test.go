package notes

import (
	"reflect"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db"
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
