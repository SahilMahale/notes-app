package server

import (
	"github.com/SahilMahale/notes-backend/internal/db"
	"github.com/SahilMahale/notes-backend/internal/notes"
	"github.com/SahilMahale/notes-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type notesService struct {
	app          *fiber.App
	DbInterface  db.DbConnection
	ip           string
	totalTickets uint
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
func (B *notesService) GetNotes(c *fiber.Ctx) error {
	return nil
}

func (B *notesService) DeleteNote(c *fiber.Ctx) error {
	return nil
}

func (B *notesService) CreateNote(c *fiber.Ctx) error {
	note := new(models.NoteRequest)

	if err := c.BodyParser(note); err != nil {
		return err
	}
	noteCtrl := notes.NewNoteController(B.DbInterface)
	noteID, err := noteCtrl.CreateBooking(note.Title, note.Body)
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
	return nil
}

func (B *notesService) StartNotesService() {
	B.initMiddleware()
	// Unauthenticated routes
	userGroup := B.app.Group("/notes")
	userGroup.Post("/create", B.CreateNote)
	userGroup.Post("/update/:noteID", B.UpdateNote)
	userGroup.Post("/delete/:noteID", B.DeleteNote)
	userGroup.Post("/get", B.GetNotes)

	/* adminGroup := B.app.Group("/admin")
	adminGroup.Post("/signup", B.CreateUser)
	adminGroup.Post("/signin", B.LoginUser) */

	B.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Booking APP Service is Running!")
	})

	// authenticated routes
	/* userGroup.Get("/info", B.GetAllUsers)
	bookingGroup := B.app.Group("/bookings")
	bookingGroup.Get("", B.GetBookings)
	bookingGroup.Post("", B.BookTickets)
	bookingGroup.Delete("/:bookID", B.DeleteBooking) */

	err := B.app.Listen(B.ip)
	if err != nil {
		log.Error(err)
		return
	}
}
