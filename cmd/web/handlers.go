package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mohit405/pkg/models"
)

// Change the signature of the home handler so it is defined as a method against application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual c
	// of r.URL.Path != "/" from this handler.

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	//use the new render helper instead
	app.render(w, r, "home.page.html", &templateData{
		Snippets: s,
	})
}

func (app *application) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	// Pat doesn't strip the colon from the named capture key, so we need to
	// get the value of ":id" from the query string instead of "id"
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.serverError(w, err)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	//use new render helper instead
	app.render(w, r, "show.page.html", &templateData{
		Snippet: s,
	})

	// data := &templateData{Snippet: s}

	// files := []string{
	// 	"./ui/html/show.page.html",
	// 	"./ui/html/base.layout.html",
	// 	"./ui/html/footer.partial.html",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) CreateSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", nil)
}

func (app *application) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	// The check of r.Method != "POST" is now superfluous and can be removed.
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	// map to hold validation errors.
	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "This field cannot be  empty"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "This field is too long (max allowed 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "This field connot be empty"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "This field cannot be empty"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "This field is invalid"
	}

	if len(errors) > 0 {
		app.render(w,r,"create.page.html",&templateData{
			FormErrors: errors,
			FormData: r.PostForm,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
