package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Create an addDefaultData helper. This takes a pointer to a tempalteData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we're not using the *http.Request parameter at the
// moment, but we will do so later.

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData{
	if td == nil{
		td = &templateData{}
	}

	td.CurrentYear= time.Now().Year()
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData){
	// Retrieve the appropriate template set from the cache based on the page 
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	
	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter
	// If there is an error, call our serverError helper and return
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil{
		app.serverError(w, err)
		return
	}
	
	// Write the contents of the buffer to the http.ResponseWriter.
	// This is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer
	buf.WriteTo(w)
}



// The serverError writes an error message and stack trace to the errorlog
// then sends a generic 500 Internal Server Error Response to the user

func (app *application) serverError(w http.ResponseWriter, err error){
	trace:= fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//Setting the frame depth to 2
	app.errorLog.Output(2, trace)
	
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}	

// The clientError helper sends a specific status code and corresponding description to the user

func (app *application) clientError(w http.ResponseWriter, status int){
	http.Error(w, http.StatusText(status), status)
}

// For consistency, we also implement a notFound helper
// This is a convenience wrapper around clientError which sends 404
func (app *application) notFound(w http.ResponseWriter){
	app.clientError(w, http.StatusNotFound)
}




