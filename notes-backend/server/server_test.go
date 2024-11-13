package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/internal/mocks"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func mockLogIn(uMock *mocks.UserOps, service *notesService, t *testing.T) string {
	login := models.UserSignin{
		Username: "dummy",
		Password: "dummy",
	}
	uMock.On("LoginUser", login.Username, login.Password).
		Return(true, helper.MyHTTPErrors{Err: nil}).Once()
	jsonLogReq, _ := json.Marshal(login)

	reqL := httptest.NewRequest("POST", "/user/signin", bytes.NewReader(jsonLogReq))
	reqL.Header.Set("Content-Type", "application/json")
	respL, errL := service.app.Test(reqL)
	assert.NoError(t, errL)
	if respL.StatusCode != fiber.StatusAccepted {
		t.Errorf("Login faied for Test: %s", t.Name())
	}
	var jwtResp models.JwtResp
	body, _ := io.ReadAll(respL.Body)
	errJ := json.Unmarshal(body, &jwtResp)
	assert.NoError(t, errJ)
	return jwtResp.Authtoken
}

func setupTestServer(t *testing.T) (*notesService, *mocks.UserOps, *mocks.NotesOps) {
	userCtrl := mocks.NewUserOps(t)
	notesCtrl := mocks.NewNotesOps(t)
	secPath, err := filepath.Abs("../secrets/")
	if err != nil {
		assert.NoError(t, err)
	}
	os.Setenv("APP_AUTH", secPath)
	service := NewNotesService("test-app", ":8001", userCtrl, notesCtrl)
	service.initMiddleware()
	userGroup := service.app.Group("/user")
	userGroup.Post("/signup", service.CreateUser)
	userGroup.Post("/signin", service.LoginUser)
	service.initAuth()
	notesGroup := service.app.Group("/notes")
	notesGroup.Post("/create", service.CreateNote)
	notesGroup.Patch("/update/:noteID", service.UpdateNote)
	notesGroup.Delete("/delete/:noteID", service.DeleteNote)
	notesGroup.Get("/get", service.GetNotes)
	notesGroup.Get("/get/:noteId", service.GetNoteByID)
	return &service, userCtrl, notesCtrl
}

func Test_notesService_CreateUser(t *testing.T) {
	service, uMock, _ := setupTestServer(t)
	tests := []struct {
		description string
		userReq     models.UserSignup
		mockMyerror helper.MyHTTPErrors
	}{
		{
			description: "successful user creation",
			userReq: models.UserSignup{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusCreated,
			},
		},
		{
			description: "failed user creation",
			userReq: models.UserSignup{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("entry already exists"),
				HttpCode: fiber.StatusBadRequest,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			uMock.On("CreateUser", test.userReq.Username, test.userReq.Email, test.userReq.Password).
				Return(test.mockMyerror).Once()

			jsonBody, _ := json.Marshal(test.userReq)
			req := httptest.NewRequest("POST", "/user/signup", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)
		})
	}
}

func Test_notesService_LoginUser(t *testing.T) {
	service, uMock, _ := setupTestServer(t)

	tests := []struct {
		loginReq    models.UserSignin
		mockMyerror helper.MyHTTPErrors
		description string
		mockIsLogIn bool
	}{
		{
			description: "LoginUser: successful",
			loginReq: models.UserSignin{
				Username: "testUser",
				Password: "testpass1234",
			},
			mockIsLogIn: true,
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusAccepted,
			},
		},
		{
			description: "LoginUser: failed",
			loginReq: models.UserSignin{
				Username: "testUser",
				Password: "testpass234",
			},
			mockIsLogIn: false,
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("password wrong"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			uMock.On("LoginUser", test.loginReq.Username, test.loginReq.Password).
				Return(test.mockIsLogIn, test.mockMyerror).Once()

			jsonBody, _ := json.Marshal(test.loginReq)
			req := httptest.NewRequest("POST", "/user/signin", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)
			if test.mockMyerror.HttpCode == fiber.StatusAccepted {
				var jwtResp models.JwtResp
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &jwtResp)
				assert.NoError(t, err)
				authToken := jwtResp.Authtoken
				fmt.Println(authToken)
			}
		})
	}
}

func TestCreateNote(t *testing.T) {
	service, uMock, nMock := setupTestServer(t)
	tests := []struct {
		description   string
		noteReq       models.NoteRequest
		mockResponse  string
		mockMyerrResp helper.MyHTTPErrors
	}{
		{
			description: "successful note creation",
			noteReq: models.NoteRequest{
				Title: "Test Note",
				Body:  "Test Body",
			},

			mockResponse: "note-123",
			mockMyerrResp: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusAccepted,
			},
		},
		{
			description: "failed note creation",
			noteReq: models.NoteRequest{
				Title: "Test Note",
				Body:  "Test Body",
			},
			mockResponse: "",
			mockMyerrResp: helper.MyHTTPErrors{
				Err:      fmt.Errorf("database error"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			nMock.On("CreateNote", test.noteReq.Title, test.noteReq.Body).
				Return(test.mockResponse, test.mockMyerrResp).Once()
			authToken := mockLogIn(uMock, service, t)
			jsonBody, _ := json.Marshal(test.noteReq)
			req := httptest.NewRequest("POST", "/notes/create", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			bearer := fmt.Sprintf("Bearer %s", authToken)
			req.Header.Set("Authorization", bearer)

			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerrResp.HttpCode, resp.StatusCode)

			if test.mockMyerrResp.HttpCode == fiber.StatusAccepted {
				var noteResp models.NoteResp
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &noteResp)
				assert.NoError(t, err)
				assert.Equal(t, test.noteReq.Title, noteResp.Title)
				assert.Equal(t, test.noteReq.Body, noteResp.Body)
				assert.Equal(t, test.mockResponse, noteResp.NoteID)
			}
		})
	}
}

func TestGetNoteByID(t *testing.T) {
	service, uMock, nMock := setupTestServer(t)

	mockNotes := models.NoteResp{
		NoteID: "note-1", Title: "Note 1", Body: "Body 1",
	}

	tests := []struct {
		description string
		noteID      string
		mockNote    models.NoteResp
		mockMyerror helper.MyHTTPErrors
	}{
		{
			description: "successful get note by ID",
			noteID:      "note-1",
			mockNote:    mockNotes,
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusOK,
			},
		},
		{
			description: "failed get  note by ID",
			noteID:      "note-1",
			mockNote:    models.NoteResp{},
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("database error"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			nMock.On("GetNote", test.noteID).Return(test.mockNote, test.mockMyerror).Once()
			authToken := mockLogIn(uMock, service, t)
			reqPath := fmt.Sprintf("/notes/get/%s", test.noteID)
			req := httptest.NewRequest("GET", reqPath, nil)
			bearer := fmt.Sprintf("Bearer %s", authToken)
			req.Header.Set("Authorization", bearer)
			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)

			if test.mockMyerror.HttpCode == fiber.StatusOK {
				var note models.NoteResp
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &note)
				assert.NoError(t, err)
				assert.Equal(t, test.mockNote.NoteID, note.NoteID)
			}
		})
	}
}

func TestGetNotes(t *testing.T) {
	service, uMock, nMock := setupTestServer(t)

	mockNotes := models.NotesResp{
		{NoteID: "note-1", Title: "Note 1", Body: "Body 1"},
		{NoteID: "note-2", Title: "Note 2", Body: "Body 2"},
	}

	tests := []struct {
		description string
		mockNotes   models.NotesResp
		mockMyerror helper.MyHTTPErrors
	}{
		{
			description: "successful get all notes",
			mockNotes:   mockNotes,
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusOK,
			},
		},
		{
			description: "failed get all notes",
			mockNotes:   nil,
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("database error"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			nMock.On("GetAllNotes").Return(test.mockNotes, test.mockMyerror).Once()
			authToken := mockLogIn(uMock, service, t)

			req := httptest.NewRequest("GET", "/notes/get", nil)
			bearer := fmt.Sprintf("Bearer %s", authToken)
			req.Header.Set("Authorization", bearer)
			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)

			if test.mockMyerror.HttpCode == fiber.StatusOK {
				var notes models.NotesResp
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &notes)
				assert.NoError(t, err)
				assert.Equal(t, len(test.mockNotes), len(notes))
				assert.Equal(t, test.mockNotes[0].NoteID, notes[0].NoteID)
			}
		})
	}
}

func TestDeleteNote(t *testing.T) {
	service, uMock, nMock := setupTestServer(t)

	tests := []struct {
		description string
		noteID      string
		mockMyerror helper.MyHTTPErrors
	}{
		{
			description: "successful note deletion",
			noteID:      "note-123",
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusOK,
			},
		},
		{
			description: "failed note deletion",
			noteID:      "note-456",
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("note not found"),
				HttpCode: fiber.StatusOK,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			nMock.On("DeleteNote", test.noteID).Return(test.mockMyerror).Once()

			authToken := mockLogIn(uMock, service, t)
			req := httptest.NewRequest("DELETE", "/notes/delete/"+test.noteID, nil)
			bearer := fmt.Sprintf("Bearer %s", authToken)
			req.Header.Set("Authorization", bearer)
			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)
		})
	}
}

func TestUpdateNote(t *testing.T) {
	service, uMock, nMock := setupTestServer(t)

	tests := []struct {
		description  string
		noteID       string
		updateReq    models.NotePatchRequest
		mockResponse models.NotePatchResp
		mockMyerror  helper.MyHTTPErrors
	}{
		{
			description: "successful note update",
			noteID:      "note-123",
			updateReq: models.NotePatchRequest{
				Title: "Updated Title",
				Body:  "Updated Body",
			},
			mockResponse: models.NotePatchResp{
				NoteID: "note-123",
				Title:  "Updated Title",
				Body:   "Updated Body",
			},
			mockMyerror: helper.MyHTTPErrors{
				Err:      nil,
				HttpCode: fiber.StatusAccepted,
			},
		},
		{
			description: "failed note update",
			noteID:      "note-456",
			updateReq: models.NotePatchRequest{
				Title: "Updated Title",
				Body:  "Updated Body",
			},
			mockResponse: models.NotePatchResp{},
			mockMyerror: helper.MyHTTPErrors{
				Err:      fmt.Errorf("note not found"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			nMock.On("UpdateNote", test.noteID, test.updateReq).
				Return(test.mockResponse, test.mockMyerror).Once()
			authToken := mockLogIn(uMock, service, t)

			jsonBody, _ := json.Marshal(test.updateReq)
			req := httptest.NewRequest("PATCH", "/notes/update/"+test.noteID, bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			bearer := fmt.Sprintf("Bearer %s", authToken)
			req.Header.Set("Authorization", bearer)

			resp, err := service.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)

			if test.mockMyerror.HttpCode == fiber.StatusAccepted {
				var noteResp models.NotePatchResp
				body, _ := io.ReadAll(resp.Body)
				err = json.Unmarshal(body, &noteResp)
				assert.NoError(t, err)
				assert.Equal(t, test.mockResponse, noteResp)
			}
		})
	}
}
