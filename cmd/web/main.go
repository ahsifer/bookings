package main

import (
	"fmt"
	"github.com/ahsifer/bookings/pkg/config"
	"github.com/ahsifer/bookings/pkg/handlers"
	"github.com/ahsifer/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const ipAddress = "0.0.0.0"

// Main function
var appConfig config.AppConfig

const port = 8080

func main() {
	appConfig.UseCache = false
	//Chane the following to true when you are using ssl in production
	appConfig.InProduction = false

	//Start Session code
	session := scs.New()
	session.Lifetime = 3 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.Secure = appConfig.InProduction
	session.Cookie.Path = "/"
	session.Cookie.SameSite = http.SameSiteLaxMode
	appConfig.Session = session
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(fmt.Sprintf("Error In Creating Template Cache : %v \n", err))
	}
	appConfig.TemplateCache = templateCache

	//Call CachePasser to make appConfig available in the render package
	render.CachePasser(&appConfig)

	//Create new repository and pass the appConfig to be available in the handlers file
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", ipAddress, port),
		Handler: routes(&appConfig),
	}
	//Start our server
	fmt.Printf("Starting server on Port: %v\n", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("unable to start the server with error:", err)
	}

}
