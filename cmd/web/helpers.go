package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

/*
Create an addDefault helper. This takes a pointer to template data struct, adds the current year to the currentYear field and then returns the pointer. Again we are not using the *http.Request parameter in this helper, but we'll need it later in the book when we add more dynamic data to the templates.
*/

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()

	return td
}

/*
	The Render function is a wrapper around the template.ExecuteTemplate() method that we talked about earlier. This means that we can call it directly in our handlers to render a template with the data from a templateData struct. The first parameter is the http.ResponseWriter where we can write the output. The second parameter is the name of the template file to render. The third parameter is a pointer to a templateData struct containing the dynamic data that we want to pass to the template. If the template doesn't exist in the cache, we call the serverError helper method to send a 500 Internal Server Error response to the user. If there's an error when rendering the template, we call the serverError helper again to send a 500 response.
*/

func (app *application) Render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	/*
		Retieve the appropriate template set from the cache based on the page number
		like (home.page.tmpl). if no entry exists in the cache with the provided name,
		call the serveError helper method that we made earilier
	*/

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Execute the template set, passing in any dynamic data
	err := ts.Execute(w, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

/*
	serveError is a helper that writes an error message and stack trace to the errorLog
	then sends a generic 500 Internal Server Error response to the user.
	debug.stack() function returns a stack trace for the goroutine that calls it.
*/

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

/*
	The clientError helper sends a specific status code and corresponding to describe
	to the user. We'll use this later in the book to send responses like 400 "Bad
	Request" when there's a problem with the request that the user sent.
*/

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

/*
	For consistency, we'll also implement a notFound helper. This is simply a
	convenience wrapper around clientError which sends a 404 Not Found response
	the user.
*/

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
