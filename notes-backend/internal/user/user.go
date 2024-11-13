package user

import (
	"fmt"

	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/helper"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserDataController struct {
	DbInterface db.DbConnection
}

type UserOps interface {
	CreateUser(username, email, pass string) helper.MyHTTPErrors
	LoginUser(username, pass string) (bool, helper.MyHTTPErrors)
}

func NewUserController(db db.DbConnection) UserDataController {
	return UserDataController{
		DbInterface: db,
	}
}

func (u UserDataController) CreateUser(username, email, pass string) helper.MyHTTPErrors {
	var user db.User
	hashPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return helper.MyHTTPErrors{
			Err:      fmt.Errorf("internal error while creating Hash"),
			HttpCode: fiber.StatusInternalServerError,
		}
	}
	user = db.User{Username: username, Email: email, Pass: string(hashPass)}

	if err := u.DbInterface.Db.Create(&user); err.Error != nil {
		myerr := helper.ErrorMatch(err.Error)
		return myerr
	}
	return helper.MyHTTPErrors{
		Err: nil,
	}
}

func (u UserDataController) LoginUser(username, pass string) (bool, helper.MyHTTPErrors) {
	user := db.User{Username: username}
	res := u.DbInterface.Db.First(&user)
	if res.Error != nil {
		myerr := helper.ErrorMatch(res.Error)
		return false, myerr
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(pass))
	if err != nil {
		myerr := helper.ErrorMatch(err)
		return false, myerr
	}
	return true, helper.MyHTTPErrors{
		Err: nil,
	}
}
