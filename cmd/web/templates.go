package main

import (
	"html/template"
	"mapcrafter2048/snippet-box/pkg/models"
	"path/filepath"
)

/*
Define a templateData type to act as the holding structure for
any dynamic data that we want to pass to our HTML templates.
At the moment it only contains one field, but we'll add more
to it as the build progresses.
*/
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as cache

	cache := map[string]*template.Template{}
	/*
		Use the filepath.Glob function to get a slice of all filepaths with
		the extension '.page.tmpl'. This essentially gives us a slice of all the
		'page' templates for the application.
	*/
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file pat
		// and assign it to the name variable.
		name := filepath.Base(page)

		// parse the page template file in a template set

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		/*
			Use the ParseGlob method to add any 'partial' templates to the
			template set (in our case, it's just the 'footer' partial at the
			moment).
		*/

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.

		cache[name] = ts

	}

	return cache, nil
}
