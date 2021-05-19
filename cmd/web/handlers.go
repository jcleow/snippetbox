package main

import (
	"fmt"
	"net/http"
	"strconv"

	// "text/template"

	"github.com/jcleow/snippetbox/pkg/models"
)

// Change the signature of the home handler so it is defined as a method against
// the *application
func (app *application) home(w http.ResponseWriter, r *http.Request) { 
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request){
	// Pat doesn't strip the colon from the named capture key, so we need to 
	// get the value of ":id" from the query string instead of "id"
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1{
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a 
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord{
		app.notFound(w)
		return
	} else if err != nil{
		app.serverError(w, err)
		return
	}

	//Use the new render helper
	app.render(w,r,"show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request){
	// The check of r.Method != "POST" is not superfluous and can be removed
	// if r.Method != "POST"{
	// 	w.Header().Set("Allow", "POST")
	// 	// http.Error(w, "Method Not Allowed", 405)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// Create some variables holding dummy data. We'll remove these later on 
	// during the build
	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji, \nBut slowly, slowly!"
	expires := "7"

	//Pass the data to the Snippet model.Insert() method, receiving the 
	// ID of the new record back.

	id, err := app.snippets.Insert(title,content,expires)
	if err != nil{
		app.serverError(w, err)
		return
	}
	// Redirect user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d",id), http.StatusSeeOther)
}

// Add a new create snippetform handler, for now returns a placeholder response
func (app * application) createSnippetForm(w http.ResponseWriter, r *http.Request){
	app.render(w, r, "create.page.tmpl",nil)
}

func (app *application) recoverPanic(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		// Create a deferred function (which will always be run in the event)
		// of a panic as Go unwinds the stack
		defer func() {
			// Use the builtin recover function to check if there has been a 
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection": close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500 
				// Internal Server Response

				// Errorf creates a new error object containing the default 
				// textual representation of the interface{} value and
				// passes this error to serverError
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)

	})
}