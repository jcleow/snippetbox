package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

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



