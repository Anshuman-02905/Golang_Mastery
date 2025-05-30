package main //package declaration

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/driver"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/handlers"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/helpers"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager //VARIABLE SHADOWING
var infolog *log.Logger
var errorlog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	defer close(app.MailChan)

	listenForMail()

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting server at Port :8080")
	//http.ListenAndServe(portNumber, nil) // nil means default ServeMux //first 1024 ports are priviledged

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//Change this to true when in production
	app.InProduction = false
	infolog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog
	errorlog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	session = scs.New()
	gob.Register(models.Reservation{})
	gob.Register(models.Restriction{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.RoomRestriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//Connect to database
	log.Println("Connecting to database ...")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=admin password=admin")
	if err != nil {
		log.Fatal("Cannot connect to database ! Dying")
		return nil, err
	}
	fmt.Println("Connected to Database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false
	//repository pattern
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
