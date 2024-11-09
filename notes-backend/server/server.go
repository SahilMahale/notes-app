package server

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/notes"
	"github.com/SahilMahale/notes-backend/internal/user"
	"github.com/SahilMahale/notes-backend/models"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

type notesService struct {
	app         *fiber.App
	DbInterface db.DbConnection
	ip          string
}

type notesServicer interface {
	GetNotes(c *fiber.Ctx) error
	DeleteNote(c *fiber.Ctx) error
	CreateNote(c *fiber.Ctx) error
	UpdateNote(c *fiber.Ctx) error
	StartNotesService(c *fiber.Ctx) error
}

func NewNotesService(appname, ip string, db db.DbConnection) notesService {
	return notesService{
		app: fiber.New(fiber.Config{
			AppName:       appname,
			StrictRouting: true,
			ServerHeader:  "NotesService",
		}),
		ip:          ip,
		DbInterface: db,
	}
}

func (B *notesService) initMiddleware() {
	// Adding logger to the app
	B.app.Use(requestid.New())
	B.app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))
	B.app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	B.app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://localhost:4000,http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}

func (B *notesService) initAuth() {
	secretsFolderPath := os.Getenv("APP_AUTH")
	if secretsFolderPath == "no-auth" || secretsFolderPath == "" {
		// run app without jwt auth
		return
	}
	privateKeyPath := fmt.Sprintf("%s/private_key.pem", secretsFolderPath)
	publicKeyPath := fmt.Sprintf("%s/public_key.pem.pub", secretsFolderPath)
	err := readPrivateKeyFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	err = readPublicKeyFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	B.app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.RS256,
			Key:    publicKey,
		},
		ContextKey: "acces-key-token",
	}))
}

func (B *notesService) GetNotes(c *fiber.Ctx) error {
	notesCtrl := notes.NewNoteController(B.DbInterface)
	notesList, err := notesCtrl.GetAllNotes()
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(notesList)
}

func (B *notesService) GetNoteByID(c *fiber.Ctx) error {
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}
	notesCtrl := notes.NewNoteController(B.DbInterface)
	note, err := notesCtrl.GetNote(noteID)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(note)
}

func (B *notesService) DeleteNote(c *fiber.Ctx) error {
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}
	notesCtrl := notes.NewNoteController(B.DbInterface)
	err := notesCtrl.DeleteNote(noteID)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Note with notedID: %s is successfully deleted.\n", noteID))
}

func (B *notesService) CreateNote(c *fiber.Ctx) error {
	note := new(models.NoteRequest)

	if err := c.BodyParser(note); err != nil {
		return err
	}
	noteCtrl := notes.NewNoteController(B.DbInterface)
	noteID, err := noteCtrl.CreateNote(note.Title, note.Body)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}

	noteResp := models.NoteResp{
		NoteID: noteID,
		Title:  note.Title,
		Body:   note.Body,
	}
	return c.Status(fiber.StatusAccepted).JSON(noteResp)
}

func (B *notesService) UpdateNote(c *fiber.Ctx) error {
	var notesCtrl notes.NotesOps
	noteID := ""
	if noteID = c.Params("noteID"); noteID == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Error: notedID not specified")
	}

	note := new(models.NotePatchRequest)

	if err := c.BodyParser(note); err != nil {
		return err
	}
	notesCtrl = notes.NewNoteController(B.DbInterface)
	noteResp, err := notesCtrl.UpdateNote(noteID, *note)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}
	return c.Status(fiber.StatusAccepted).JSON(noteResp)
}

func (B *notesService) CreateUser(c *fiber.Ctx) error {
	var userCtrl user.UserOps
	u := new(models.UserSignup)

	if err := c.BodyParser(u); err != nil {
		return err
	}

	userCtrl = user.NewUserController(B.DbInterface)

	errP := userCtrl.CreateUser(u.Username, u.Email, u.Password)
	if errP.Err != nil {
		return c.Status(errP.HttpCode).SendString(errP.Err.Error())
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (B *notesService) LoginUser(c *fiber.Ctx) error {
	var userCtrl user.UserOps
	u := new(models.UserSignin)
	if err := c.BodyParser(u); err != nil {
		return err
	}
	userCtrl = user.NewUserController(B.DbInterface)

	_, err := userCtrl.LoginUser(u.Username, u.Password)
	if err.Err != nil {
		return c.Status(err.HttpCode).SendString(err.Err.Error())
	}

	// Create a token based on user
	atoken, errp := makeTokenWithClaims(u.Username)

	if errp != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(errp.Error())
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"auth_token": atoken})
}

func (B *notesService) StartNotesService() {
	B.initMiddleware()
	// Unauthenticated routes
	userGroup := B.app.Group("/user")
	userGroup.Post("/signup", B.CreateUser)
	userGroup.Post("/signin", B.LoginUser)

	B.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Booking APP Service is Running!")
	})
	B.initAuth()
	// authenticated routes

	notesGroup := B.app.Group("/notes")
	notesGroup.Post("/create", B.CreateNote)
	notesGroup.Patch("/update/:noteID", B.UpdateNote)
	notesGroup.Delete("/delete/:noteID", B.DeleteNote)
	notesGroup.Get("/get", B.GetNotes)
	notesGroup.Get("/get/:noteId", B.GetNoteByID)

	err := B.app.Listen(B.ip)
	if err != nil {
		log.Error(err)
		return
	}
}
