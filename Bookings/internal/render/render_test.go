package render

import (
	"net/http"
	"testing"

	"github.com/Anshuman-02905/Golang_Mastery/Bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()

	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in sesssio")
	}
}

func Test_RederTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error("Session error", err)
	}

	ww := myWriter{}

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})

	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "lol.page.tmpl", &models.TemplateData{})

	if err == nil {
		t.Error("Rendered template that does not exists")
	}
}

func test_NewTemplates(t *testing.T) {
	NewRenderer(app)
}
func Test_CreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil

}
