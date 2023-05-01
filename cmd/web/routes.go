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
	dynamicMiddleware := alice.New(app.session.Enable, app.authenticate)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// Add the requireAuthenticatedUser middleware to the chain.
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.CreateSnippetForm))
	// Add the requireAuthenticatedUser middleware to the chain.
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.CreateSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.ShowSnippet))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForn))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser))

	filerServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", filerServer))

	return standardMiddleware.Then(mux)
	// above statement is equvivanlent to myMiddleware1(myMiddleware2(myMiddleware3(myHandler)))
}
