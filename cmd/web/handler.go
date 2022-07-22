package main

import (
	"fmt"
	"github.com/Onolax/SnipB/pkg/models"
	"net/http"
	"strconv"
)

//notice we have linked the function with *application method in main file
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//Path extractor
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

	//data := &templateData{Snippets: s}

	////byte slice containing all the required files so that we can easily access them
	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	////Template file reading
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//}
	//
	//// Executing the template file
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	//Extracting the query and getting the key
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

	//storage of templateData struct to hold the snippet data
	//data := &templateData{Snippet: s}
	//
	////slice of html files
	//files := []string{
	//	"./ui/html/show.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	////Parse the files
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//}
	////executing them using snippet data
	//err = ts.Execute(w, data)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	//restricting the handler to only POST requests
	//What are POST requests?
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
