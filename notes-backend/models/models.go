package models

type NoteRequest struct {
	Title string `json:"title" xml:"title" form:"title"`
	Body  string `json:"body" xml:"body" form:"body"`
}
type NotesRequest []NoteRequest

type NoteResp struct {
	NoteID string `json:"noteID" xml:"noteID" form:"noteID"`
	Title  string `json:"title" xml:"title" form:"title"`
	Body   string `json:"body" xml:"body" form:"body"`
}

type NotesResp []NoteResp
