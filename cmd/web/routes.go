package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our appliation receives
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	
	//Use the mux.Handle() fn to register the fileServer as the handler
	// for all URL paths that start with "/static/". For matching paths, we strip
	// the "/static" prefix before the request reaches the file server.
	//https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files
	// mux := http.NewServeMux()
	// mux.HandleFunc("/",app.home)
	// mux.HandleFunc("/snippet",app.showSnippet)
	// mux.HandleFunc("/snippet/create",app.createSnippet)

	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware
	// because secureHeaders is just a function, and the function returns a
	// http.Handler, we don't need to do anything else

	// wrap the existing chain with th logRequest middle ware
	// wrap the existing chain with the recoverPanic middleware

	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standardMiddleware.Then(mux)
}