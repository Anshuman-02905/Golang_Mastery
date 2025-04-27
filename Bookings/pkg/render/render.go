package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/models"
)

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// Renders tepmlate using html template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	//get the templaet form appConfig
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could no get the Template from AppCOnfig ")

	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

}
func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)
	myCache := map[string]*template.Template{}

	//get all of the files
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		log.Println("Error at  Getting all files")
		return myCache, err
	}
	//range through all files endig  with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)

		if err != nil {
			log.Printf("Error at  parsing %s file ", name)

			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			log.Printf("Error at  parsing %s file layout ", name)

			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				log.Printf("Error at  parsing %s file matches ", name)

				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil

}

//var tc = make(map[string]*template.Template)

// func RenderTemplateTest(w http.ResponseWriter, t string) {
// 	//store a parsed template in a dta structure
// 	var tmpl *template.Template
// 	var err error
// 	//check if already have the emplate

// 	_, inMap := tc[t]
// 	if !inMap {
// 		//need to vreate the template
// 		log.Println("Creating template and adding to cache")
// 		err = CreateTemplateCache(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		// we have the template in the cache
// 		log.Println("usig cacheced template")
// 	}
// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}

// }
// func CreateTemplateCache(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t), "./templates/base.layout.tmpl",
// 	}

// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err

// 	}
// 	tc[t] = tmpl
// 	return nil
// }
