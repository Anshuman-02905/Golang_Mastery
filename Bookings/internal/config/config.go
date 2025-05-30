package config

import (
	"html/template"
	"log"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

//Should not use stuff from other parts of library but can be used everywhere else

// AppConfig holds the applicatiion Config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}

//var App *AppConfig
