package main

import (
	"net/http"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v4"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(WriteToConsole)
	mux.Use(NoServer)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	return mux
}
