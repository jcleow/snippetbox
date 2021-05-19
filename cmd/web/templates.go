package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/jcleow/snippetbox/pkg/models"
)

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This essentially is a string-keyed map which acts a lookup between the names
// of a custom template functions and the functions themselves
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates
// At the moment, it only contains one field, but we'll add more to it
// as the build progresses.

type templateData struct{
	CurrentYear int
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error){
	//Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with 
	// the extension '.page.tmpl'. This essentially gives of all the 'page' templates
	// for the application

	//https://golang.org/pkg/path/filepath/#Glob --> takes in (pattern string)
	// --> returns an array of string matches and error

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil{
		return nil, err
	}

	// Loop through the pages one by one
	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path'
		// and assign it to the name variable.
		// https://golang.org/pkg/path/filepath/#Base
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		
		// The template.FuncMap must be registered with the template set before
		// they call the ParseFiles() method. This means we have to use template.New
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		// https://golang.org/pkg/text/template/#Template.Funcs
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil{
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set ( in our case, it is just the 'base' layout at the moment)
		// https://golang.org/pkg/html/template/#ParseGlob
		ts, err = ts.ParseGlob(filepath.Join(dir,"*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the 
		// template set (in our case, it's just the 'footer' partial at the moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}
	// Return the map.
	return cache, nil

}