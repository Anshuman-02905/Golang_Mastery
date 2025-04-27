package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

//Should not use stuff from other parts of library but can be used everywhere else

// AppConfig holds the applicatiion Config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InProduction  bool
	Session       *scs.SessionManager
}

//var App *AppConfig
