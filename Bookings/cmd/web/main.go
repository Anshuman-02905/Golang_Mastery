package main //package declaration

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/handlers"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager //VARIABLE SHADOWING

func main() {
	//Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session=session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	//repository pattern
	repo := handlers.NewRepo((&app))
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println("Starting server at Port :8080")
	//http.ListenAndServe(portNumber, nil) // nil means default ServeMux //first 1024 ports are priviledged

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
}
