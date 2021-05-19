package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jcleow/snippetbox/pkg/models/mysql"

	// since we don't use anything in the mysql package
	// we try to import it normally
	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies
type application struct{
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Defines a new CLI flag with name addr, default value
	// some short help text explaining what the flag controls
	// value of flag will be stored in the addr variable at runtime
	addr := flag.String("addr",":4000","HTTP network address")
	dsn := flag.String("dsn","web:Jityong1@/snippetbox?parseTime=true", "MySQL database connection")
	
	// use flag.Parse() function to parse the command-line
	// reads in the command-line flag value and assigns it to the addr variable
	// Need to call this before we use the addr variable, if not it will always contain the default value of :400
	flag.Parse()

	// custom loggers created are concurrency safe
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	// To keep the main() function tidy, we put the code for creating a connection
	// pool into the separate openDB function below
	// We pass openDB() the DSN variable from the command-line flag
	db, err := openDB(*dsn)
	if err != nil{
		errorLog.Fatal(err)
	}

	//We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits
	

	// Defer is used to ensure that a function call is performed later in a
	// program’s execution, usually for purposes of cleanup
	// https://gobyexample.com/defer
	defer db.Close()

	//Initialize a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil{
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application contianing dependencies
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets :&mysql.SnippetModel{DB:db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}



	// The value returned from flag.String() function is a pointer to the flag's
	// value, not the value itself. So we need to dereference the pointer
	// log.Printf("Starting server on %s", *addr)

	infoLog.Printf("Starting server %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	serverErr := srv.ListenAndServe()
	errorLog.Fatal(serverErr)
}

func openDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql",dsn)
	if err != nil{
		return nil, err
	}
	// db.Ping() is used to establish a connection and check for any errors
	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}