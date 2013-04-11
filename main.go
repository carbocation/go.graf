package main

import (
	"database/sql"
	"net/http"
	"runtime"

	"bitbucket.org/tebeka/nrsc"
	"github.com/carbocation/forum.git/forum"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("f2LdNYi5fvo8YNdMDvI9Ggnv2OUaRiIEXFUru+v23ZxskQ"))
var router *mux.Router

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Initialize the DB in the main function so we'll have a pool of connections maintained
	db = initdb()
	defer db.Close()

	//Initialize the forum package
	forum.Initialize(db)

	//Bundled static assets are handled by nrsc
	nrsc.Handle("/static/")

	//Dynamic content is managed by handlers pointed at by the router 
	router = mux.NewRouter()

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", handler(indexHandler)).Name("index")
	g.Handle("/thread/{id:[0-9]+}", handler(threadHandler)).Name("thread")

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", handler(newThreadHandler)).Name("createThread")
	p.Handle("/login/{id:[0-9]+}", handler(loginHandler)).Name("login")

	//Notify the http package about our router
	http.Handle("/", router)

	//Launch the server
	if err := http.ListenAndServe("localhost:9999", nil); err != nil {
		panic(err)
	}
}

func initdb() *sql.DB {
	db, err := sql.Open("postgres", "dbname=forumtest sslmode=disable")
	if err != nil {
		panic(err)
	}

	return db
}
