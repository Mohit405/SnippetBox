package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

// we moved routes to this file to avoid main file overcrowding..
// these are routes are now standalone from the main package ..
func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)

	// mux.HandleFunc("/snippet", app.ShowSnippet)
	// mux.HandleFunc("/snippet/create", app.CreateSnippet)

	// fileserver := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.CreateSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.CreateSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.ShowSnippet))

	filerServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", filerServer))

	return standardMiddleware.Then(mux)
	// above statement is equvivanlent to myMiddleware1(myMiddleware2(myMiddleware3(myHandler)))
}
