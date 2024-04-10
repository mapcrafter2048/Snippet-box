package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"mapcrafter2048/snippet-box/pkg/models/mysql"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	/*
		Define a new command-line flag with the name 'addr', a default value of
		and some short help text explaining what the flag controls. The value of
		flag will be stored in the addr variable at runtime.
		The flag.String() function returns a string pointer (string) so we need
	*/

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	/*
		Importantly, we use the flag.Parse() function to parse the command-line
		This reads in the command-line flag value and assigns it to the addr
		variable.
	*/

	flag.Parse()

	/*
		Use log.New() to create a logger for writing information messages. This
		three parameters: the destination to write the logs to (os.Stdout), a st
		prefix for message (INFO followed by a tab), and flags to indicate what
		additional information to include (local date and time). Note that the fl
		are joined using the bitwise OR operator |.
	*/

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	/*
		Create a logger for writing error messages in the same way, but use stde
		the destination and use the log.Lshortfile flag to include the relevant
		file name and line number.
	*/

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.LUTC|log.Llongfile)

	/*
		To keep the main() function tidy I've put the code for creating a connec
		pool into the separate openDB() function below. We pass openDB() the DSN
		from the command-line flag
	*/

	db, err := openDB(*dsn)

	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	/*
		Use the http.NewServeMux() function to initialize a new servemux
		mux := http.NewServeMux()

		Use the mux.HandleFunc() method to register a new handler function for the "/"
		mux.HandleFunc("/", app.Home)

		Use the mux.HandleFunc() method to register a new handler function for the "/snippet" URL pattern
		mux.HandleFunc("/snippet", app.ShowSnippet)

		Use the mux.HandleFunc() method to register a new handler function for the "/snippet/create" URL pattern
		mux.HandleFunc("/snippet/create", app.CreateSnippet)

		Create a file server which serves files out of the "./ui/static" directo
		Note that the path given to the http.Dir function is relative to the pro
		directory root.
		fileServer := http.FileServer(http.Dir("./ui/static/"))

		Use the mux.Handle() function to register the file server as the handler
		all URL paths that start with "/static/". For matching paths, we strip t
		"/static" prefix before the request reaches the file server.

		mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	*/

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.Routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)

	// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
	// for a given DSN.

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
