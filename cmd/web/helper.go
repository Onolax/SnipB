package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

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
