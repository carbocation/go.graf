package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"bitbucket.org/tebeka/nrsc"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var db *sql.DB
var store = sessions.NewCookieStore([]byte("something-very-secret"))
var router *mux.Router

func main() {
	// Initialize the DB in the main function so we'll have a pool of connections maintained
	db = initdb()
	defer db.Close()
	
	//Bundled static assets are handled by nrsc
	nrsc.Handle("/static/")

	//Initialize our router
	router = mux.NewRouter()

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", handler(indexHandler)).Name("index")
	g.Handle("/thread/{id:[0-9]+}", handler(threadHandler)).Name("thread")
	g.Handle("/css/{file}", handler(cssHandler))
	//g.HandleFunc("/static/", staticHandler)

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

func mainCss() string {
	return `
div .comment {
	padding-left: 100px;
}
`
}

func PrintNestedComments(el *binarytree.Tree, ch chan string) {
	if el == nil {
		return
	}

	ch <- "<div class=\"comment\">"

	//Self
	e := el.Value.(forum.Entry)
	ch <- fmt.Sprintf("Title: %s", e.Title)

	//Children are embedded
	PrintNestedComments(el.Left(), ch)
	ch <- "</div>"

	//Siblings are parallel
	PrintNestedComments(el.Right(), ch)
}
