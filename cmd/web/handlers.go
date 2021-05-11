package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Change the signature of the home handler so it is defined as a method against
// the *application
func (app *application) home(w http.ResponseWriter, r *http.Request) { 
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// Initialize a slice containing the paths to the two files
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Must be absolute or relative file path
	// ts refers to template set
	ts, err := template.ParseFiles(files...)
	if err != nil {
		// log.Println(err.Error())
		// app.errorLog.Println(err.Error())
		// http.Error(w,"Internal Server Error", 500)
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w,nil)
	if err != nil{
		// log.Println(err.Error())
		// app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1{
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request){
	if r.Method != "POST"{
		w.Header().Set("Allow", "POST")
		// http.Error(w, "Method Not Allowed", 405)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}