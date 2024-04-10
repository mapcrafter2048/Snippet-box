package main

import (
	"fmt"
	"html/template"
	"mapcrafter2048/snippet-box/pkg/models"
	"net/http"
	"strconv"
)

// define a home handler which writes a byte slice containing "Hello from snippet" as a response body
func (app *application) Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.

	// s, err := app.snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }

	s, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// // We then use the Execute() method on the template set to write the templa
	// // content as the response body. The last parameter to Execute() represents
	// // dynamic data that we want to pass in, which for now we'll leave as nil

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

	// w.Write([]byte("Hello from snippet"))
}

// Add a showsnippet handler function
func (app *application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() function. If it can't
	// be converted to an integer, or the value is less than 1, we return a 404
	// not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
	}

	data := &templateData{Snippet: s}

	files := []string{
		"./ui/html/show.page.tmpl",
		// "./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partial that we made earlier.

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Pass in the templateData Struct when executing the template set
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}

	// fmt.Fprintf(w, "%v", s)
	// // w.Write([]byte("Display a specific snippet "))
	// fmt.Fprintf(w, "Display a specific snippet with ID %d", id)

}

// Add a createSnippet handler function
func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		// Use the http.Error() function to send a 405 status code and "Method N
		// Allowed" string as the response body
		w.Header().Set("Allow", "POST") // sets the allow header to POST so that the client knows that the only allowed method is POST and other methods are not allowed
		/*
			other header methods incllude:
			w.Header().Set("Content-Type", "application/json") // sets the content type header to application/json
			w.Header().add("Content-Type", "application/json") // adds a new value to the content type header
			w.Header().Del("Content-Type") // deletes the content type header
			w.Header().Get("Content-Type") // gets the value of the content type header
		*/
		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d ", id), http.StatusSeeOther)

	// w.Write([]byte("Create a new snippet")) //sends 200 OK so to send a non 200 Ok status code we must call w.WriterHeader(http.StatusOK) before writing the response body
}
