package main

import (
	"fmt"
	"html/template"
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

	//byte slice containing all the required files so that we can easily access them
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	//Template file reading
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
	}

	// Executing the template file
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	//Extracting the query and getting the key
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	//restricting the handler to only POST requests
	//What are POST requests?
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Create a new snippet..."))
}
