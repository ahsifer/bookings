package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ahsifer/bookings/pkg/config"
	"github.com/ahsifer/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var appConfig config.AppConfig

func CachePasser(tc *config.AppConfig) {
	appConfig = *tc
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	//here you can add the default data that might be used by all or almost all templates to the passed template data
	return td
}

// TemplateRender this function is used to render html pages
func TemplateRender(w http.ResponseWriter, templateName string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if appConfig.UseCache {
		tc = appConfig.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	currentTemplate, ok := tc[templateName]
	if !ok {
		log.Fatal(fmt.Sprintf("No Such Template in Template Cache: %v \n", templateName))
	}
	buffer := new(bytes.Buffer)
	td = AddDefaultData(td)
	err := currentTemplate.Execute(buffer, td)
	if err != nil {
		fmt.Printf("Cannot execute the template to the buffer: %v \n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = buffer.WriteTo(w)
	if err != nil {
		fmt.Printf("Cannot write data from the buffer to the ResponseWriter %v \n", err)
	}
}

var functions = template.FuncMap{}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if len(pages) == 0 || err != nil {
		return myCache, errors.New("malformed pattern or no file exists")
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
