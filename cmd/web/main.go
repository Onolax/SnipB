package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	//used a flag and stored in value in addr and parsed it so that it can be altered during runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	//made two different logger for different situations
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//made the pointer so that the router can use handler with application methods
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// used the http.Server so that it uses our new error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		//used the routes method from routes.go
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}
