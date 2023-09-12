package handlers

import (
	"fmt"
	"github.com/ahsifer/bookings/pkg/config"
	"github.com/ahsifer/bookings/pkg/models"
	"github.com/ahsifer/bookings/pkg/render"

	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepo creates new Repository
func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

// NewHandlers sets the repository for the new handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home page requests handler
func (r *Repository) Home(response http.ResponseWriter, request *http.Request) {
	remoteIP := request.RemoteAddr
	r.App.Session.Put(request.Context(), "My_IP", remoteIP)
	render.TemplateRender(response, "home.page.tmpl", &models.TemplateData{})
}

// About page requests handler
func (r *Repository) About(response http.ResponseWriter, request *http.Request) {
	remoteIP := r.App.Session.GetString(request.Context(), "My_IP")
	AboutStringMap := map[string]string{
		"test": fmt.Sprintf("This is test string passed to the about page and your IP addrress is: %v", remoteIP),
	}
	render.TemplateRender(response, "about.page.tmpl", &models.TemplateData{
		StringMap: AboutStringMap,
	})
}
