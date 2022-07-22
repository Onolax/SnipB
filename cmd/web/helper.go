package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

//handler to send server side side error which sends 500 status error
func (app *application) serverError(w http.ResponseWriter, err error) {

	//debug.Stack() just traces the execution path of current goroutines
	//which is helpful when debugging the code
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	//Output prints till the given depth of trace
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//handler to handle client side error(mostly error like 400)
func (app *application) clientError(w http.ResponseWriter, status int) {

	//http.StatusText simply generates human-friendly
	//text representation of status code
	http.Error(w, http.StatusText(status), status)
}

//this is helper of clientError which simply sends 404 Status code
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	//making a buffer
	buf := new(bytes.Buffer)

	//writing in buffer before screen to find any errors
	//and not print half completed response
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	//writing from buffer to response writer
	buf.WriteTo(w)

}
