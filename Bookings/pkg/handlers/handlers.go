package handlers

import (
	"net/http"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/config"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/models"
	"github.com/Anshuman-02905/Golang_Mastery/Bookings/pkg/render"
)

// Template data hold data  sent from handlers to Templates

// REpo the repositody used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct { //Repository patter

	App *config.AppConfig
}

// Creates the NewRpositoyr
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// New Handlers sers the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr

	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some buisness logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	//send some data
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

}
