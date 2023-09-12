package main

import (
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit The Page")
		next.ServeHTTP(w, r)
	})

}
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return next
}

// SessionLoad function is used to load the session from the user and track users cookie feel free to see what LoadAndSave source code
func SessionLoad(next http.Handler) http.Handler {
	return appConfig.Session.LoadAndSave(next)
}
