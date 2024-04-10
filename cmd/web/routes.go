package main

import "net/http"

/*
Routes Define a new application struct. This holds the application-wide dependencies
for the web application. For now we'll only include fields for the two custom
loggers, but we'll add more to it as the build progresses.
*/
func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet", app.ShowSnippet)
	mux.HandleFunc("/snippet/create", app.CreateSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
