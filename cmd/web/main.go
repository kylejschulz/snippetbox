package main

import (
    "database/sql"
    "flag"
    "log"
		"html/template"
    "net/http"
	  "os"
		"snippetbox.kyleschulz.net/internal/models"
		"github.com/go-playground/form/v4" // New import
    _ "github.com/go-sql-driver/mysql"
)

type application struct {
    errorLog *log.Logger
    infoLog * log.Logger
		snippets *models.SnippetModel
		templateCache map[string]*template.Template
		formDecoder *form.Decoder
}

func main() {
    addr := flag.String("addr", ":4000", "HTTP network address")

		dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

    flag.Parse()

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

		// To keep the main() function tidy I've put the code for creating a connection
		// pool into the separate openDB() function below. We pass openDB() the DSN
		// from the command-line flag.
		db, err := openDB(*dsn)
		if err != nil {
				errorLog.Fatal(err)
		}

		// We also defer a call to db.Close(), so that the connection pool is closed
		// before the main() function exits.
		defer db.Close()

		// Initialize a new template cache...
		templateCache, err := newTemplateCache()
		if err != nil {
			errorLog.Fatal(err)
		}
		// Initialize a decoder instance...
		formDecoder := form.NewDecoder()

		// And add it to the application dependencies.
		app := &application{
			errorLog: errorLog,
			infoLog: infoLog,
			snippets: &models.SnippetModel{DB: db},
			templateCache: templateCache,
			formDecoder: formDecoder,
		}

		srv := &http.Server{
        Addr: *addr,
        ErrorLog: errorLog,
        Handler: app.routes(),
    }

    infoLog.Printf("Starting server on %s", *addr)
    err = srv.ListenAndServe()
    errorLog.Fatal(err)
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
