package notes

import (
	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/google/uuid"
)

type NotesController struct {
	DbInterface db.DbConnection
}

type NotesOps interface {
	CreateNote(title, body string) (string, helper.MyHTTPErrors)
	UpdateNote(noteId, title, body string) (string, helper.MyHTTPErrors)
	DeleteNote(noteId string) helper.MyHTTPErrors
	GetAllNotes() (models.NotesResp, helper.MyHTTPErrors)
	GetNote(noteID string) (models.NoteResp, helper.MyHTTPErrors)
}

func NewNoteController(db db.DbConnection) NotesController {
	return NotesController{
		DbInterface: db,
	}
}

func (b NotesController) CreateBooking(title, body string) (string, helper.MyHTTPErrors) {
	noteId := uuid.NewString()
	note := db.Note{
		NoteID: noteId,
		Title:  title,
		Body:   body,
	}

	if err := b.DbInterface.Db.Create(&note); err.Error != nil {
		return "", helper.ErrorMatch(err.Error)
	}
	return noteId, helper.MyHTTPErrors{
		Err: nil,
	}
}

func (n NotesController) GetAllNotes() (models.NotesResp, helper.MyHTTPErrors) {

	notesList := models.NotesResp{}
	notes := []db.Note{}
	res := n.DbInterface.Db.Find(&notes)
	if res.Error != nil {
		myerr := helper.ErrorMatch(res.Error)
		return nil, myerr
	}
	for _, note := range notes {
		notesList = append(notesList, models.NoteResp{
			NoteID: note.NoteID,
			Title:  note.Title,
			Body:   note.Body,
		})
	}
	return notesList, helper.MyHTTPErrors{
		Err: nil,
	}
}
