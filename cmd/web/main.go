package main

import (
	"database/sql"
	"flag"
	"github.com/Onolax/SnipB/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	//used a flag and stored in value in addr and parsed it so that it can be altered during runtime
	addr := flag.String("addr", ":4000", "HTTP network address")
	//flag for MySQL
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data")
	flag.Parse()

	//made two different logger for different situations
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	//made the pointer so that the router can use handler with application methods
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// used the http.Server so that it uses our new error logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		//used the routes method from routes.go
		Handler: app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
