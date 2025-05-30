package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig
var pathToTemplates = "./templates"
var functions = template.FuncMap{}

func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(),"user_id"){
		td.IsAuthenticated=1
	}

	return td
}

//  tepmlate using html template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	//get the templaet form appConfig
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		//log.Fatal("Could no get the Template from AppCOnfig ")
		return errors.New("Cannot get Template from Cache")

	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
func CreateTemplateCache() (map[string]*template.Template, error) {
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
