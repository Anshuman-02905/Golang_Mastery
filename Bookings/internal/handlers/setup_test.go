package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"

	"github.com/alexedwards/scs/v2"
)

var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}

func NoServer(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

var pathToTemplates = "./../../templates"

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		log.Println("Error at  Getting all files")
		return myCache, err
	}
	//range through all files endig  with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			log.Printf("Error at  parsing %s file ", name)

			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			log.Printf("Error at  parsing %s file layout ", name)

			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				log.Printf("Error at  parsing %s file matches ", name)

				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil

}

// SEssionLoad loads and Save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func TestMain(m *testing.M) {

	app.InProduction = false

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infolog
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorlog

	session = scs.New()
	gob.Register(models.Reservation{})
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)
	listenForMail()

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	//repository pattern
	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)
	os.Exit(m.Run())

}

func getRoutes() http.Handler {

	mux := chi.NewRouter()
	//mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	//mux.Use(NoServer)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quaters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)

	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}

	}()
}
