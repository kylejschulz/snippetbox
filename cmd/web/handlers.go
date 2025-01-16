package main

import (
   	"errors"
		"fmt"
    "html/template"
    "net/http"
    "strconv"

		"snippetbox.kyleschulz.net/internal/models"
)


func (app *application) home(w http.ResponseWriter, r *http.Request ) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}



  files := []string{
    "./ui/html/base.tmpl",
    "./ui/html/pages/home.tmpl",
    "./ui/html/partials/nav.tmpl",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    app.serverError(w, err)
    return
  }

	// Create an instance of a templateData struct holding the slice of
	// snippets.
	data := &templateData{
		Snippets: snippets,
	}

  err = ts.ExecuteTemplate(w,"base", data)
  if err != nil {
    app.serverError(w, err)
  }
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
				app.notFound(w)
				return
		}

		snippet, err := app.snippets.Get(id)
		if err != nil {
				if errors.Is(err, models.ErrNoRecord) {
						app.notFound(w)
				} else {
						app.serverError(w, err)
				}
				return
		}

		files := []string{
				"./ui/html/base.tmpl",
				"./ui/html/partials/nav.tmpl",
				"./ui/html/pages/view.tmpl",
		}

		// Parse the template files...
		ts, err := template.ParseFiles(files...)
		if err != nil {
				app.serverError(w, err)
		return
		}

		// Create an instance of a templateData struct holding the snippet data.
		data := &templateData{
				Snippet: snippet,
		}

		// And then execute them. Notice how we are passing in the snippet
		// data (a models.Snippet struct) as the final parameter?
		err = ts.ExecuteTemplate(w, "base", data)
		if err != nil {
				app.serverError(w, err)
		}
}


func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request ) {

		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			app.clientError(w, http.StatusMethodNotAllowed)
			return
		}

		// Create some variables holding dummy data. We'll remove these later on
		// during the build.
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
		expires := 7

		// Pass the data to the SnippetModel.Insert() method, receiving the
		// ID of the new record back.
		id, err := app.snippets.Insert(title, content, expires)
		if err != nil {
			app.serverError(w, err)
			return
		}

		// Redirect the user to the relevant page for the snippet.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}


