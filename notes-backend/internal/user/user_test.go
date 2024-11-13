package user

import (
	"testing"

	"github.com/SahilMahale/notes-backend/internal/db/mock"
	"github.com/SahilMahale/notes-backend/internal/helper"
)

func TestUserDataController_CreateUser(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	uctrl := NewUserController(dbMock)
	type args struct {
		username string
		email    string
		pass     string
	}
	tests := []struct {
		name string
		args args
		want helper.MyHTTPErrors
	}{
		{
			name: "CreateUser: Positive",
			args: args{
				username: "testUser",
				email:    "testUser@mail.com",
				pass:     "test1234567",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uctrl.CreateUser(tt.args.username, tt.args.email, tt.args.pass); got.Err != tt.want.Err {
				t.Errorf("UserDataController.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserDataController_LoginUser(t *testing.T) {
	dbMock := mock.NewTestDB(t)
	uctrl := NewUserController(dbMock)
	err := uctrl.CreateUser("testUser", "testUser@mail.com", "test1234567")
	if err.Err != nil {
		t.Error(err.Err)
	}

	type args struct {
		username string
		pass     string
	}
	tests := []struct {
		wantMyErr helper.MyHTTPErrors
		args      args
		name      string
		want      bool
	}{
		{
			name: "LoginUser: Positive",
			args: args{
				username: "testUser",
				pass:     "test1234567",
			},
			want: true,
			wantMyErr: helper.MyHTTPErrors{
				Err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := uctrl.LoginUser(tt.args.username, tt.args.pass)
			if got != tt.want {
				t.Errorf("UserDataController.LoginUser() got = %v, want %v", got, tt.want)
			}
			if gotErr != tt.wantMyErr {
				t.Errorf("UserDataController.LoginUser() gotErr = %v, want %v", gotErr, tt.wantMyErr)
			}
		})
	}
}
