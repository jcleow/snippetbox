package main

import "net/http"

func (app *application) routes() *http.ServeMux{

	//Use the mux.Handle() fn to register the fileServer as the handler
	// for all URL paths that start with "/static/". For matching paths, we strip
	// the "/static" prefix before the request reaches the file server.
	//https://stackoverflow.com/questions/27945310/why-do-i-need-to-use-http-stripprefix-to-access-my-static-files

	mux := http.NewServeMux()
	mux.HandleFunc("/",app.home)
	mux.HandleFunc("/snippet",app.showSnippet)
	mux.HandleFunc("/snippet/create",app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	return mux
}