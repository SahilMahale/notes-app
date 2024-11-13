package helper

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestValidateUserInfo(t *testing.T) {
	type args struct {
		username string
		email    string
	}
	tests := []struct {
		name                string
		args                args
		wantIsvalidUsername bool
		wantIsvalidEmail    bool
	}{
		{
			name: "ValidateUserInfo: Valid_UserName",
			args: args{
				username: "testUser",
				email:    "testUser@mail.com",
			},
			wantIsvalidUsername: true,
			wantIsvalidEmail:    true,
		},
		{
			name: "ValidateUserInfo: Invalid_UserInfo",
			args: args{
				username: "T",
				email:    "Tmail.com",
			},
			wantIsvalidUsername: false,
			wantIsvalidEmail:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsvalidUsername, gotIsvalidEmail := ValidateUserInfo(tt.args.username, tt.args.email)
			if gotIsvalidUsername != tt.wantIsvalidUsername {
				t.Errorf("ValidateUserInfo() gotIsvalidUsername = %v, want %v", gotIsvalidUsername, tt.wantIsvalidUsername)
			}
			if gotIsvalidEmail != tt.wantIsvalidEmail {
				t.Errorf("ValidateUserInfo() gotIsvalidEmail = %v, want %v", gotIsvalidEmail, tt.wantIsvalidEmail)
			}
		})
	}
}

func TestErrorMatch(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want MyHTTPErrors
	}{
		{
			name: "ErrorMatch: ErrDupKey",
			args: args{
				err: gorm.ErrDuplicatedKey,
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("entry already exists"),
				HttpCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "ErrorMatch: ErrRecNF",
			args: args{
				err: gorm.ErrRecordNotFound,
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("entry doesn't exist"),
				HttpCode: fiber.StatusForbidden,
			},
		},
		{
			name: "ErrorMatch: DupEntry_User",
			args: args{
				err: fmt.Errorf("Duplicate entry for key 'users.PRIMARY'"),
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("username already exists"),
				HttpCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "ErrorMatch: DupEntry",
			args: args{
				err: fmt.Errorf("Duplicate entry username already exists"),
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("entry already exists"),
				HttpCode: fiber.StatusBadRequest,
			},
		},
		{
			name: "ErrorMatch: PassWrong",
			args: args{
				err: bcrypt.ErrMismatchedHashAndPassword,
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("password wrong"),
				HttpCode: fiber.StatusForbidden,
			},
		},
		{
			name: "ErrorMatch: ForeignKey_userNotFound",
			args: args{
				err: fmt.Errorf("FOREIGN KEY"),
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("user not found"),
				HttpCode: fiber.StatusForbidden,
			},
		},
		{
			name: "ErrorMatch: Unmatched_Error",
			args: args{
				err: fmt.Errorf("Unmatched_Error"),
			},
			want: MyHTTPErrors{
				Err:      fmt.Errorf("internal server error"),
				HttpCode: fiber.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorMatch(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ErrorMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
