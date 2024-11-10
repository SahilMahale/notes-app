package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/SahilMahale/notes-backend/internal/mocks"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) (*notesService, *mocks.UserOps, *mocks.NotesOps) {
	userCtrl := mocks.NewUserOps(t)
	notesCtrl := mocks.NewNotesOps(t)
	service := NewNotesService("test-app", ":8001", userCtrl, notesCtrl)
	userGroup := service.app.Group("/user")
	userGroup.Post("/signup", service.CreateUser)
	userGroup.Post("/signin", service.LoginUser)
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
			fmt.Println("wtf", resp)
			assert.Equal(t, test.mockMyerror.HttpCode, resp.StatusCode)
		})
	}
}
