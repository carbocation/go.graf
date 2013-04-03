package main

import (
	"fmt"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"net/http"
	//"html/template"
	"database/sql"
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

	//Initialize our router
	router = mux.NewRouter()

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", handler(defaultHandler))
	g.Handle("/thread/{id:[0-9]+}", handler(threadHandler))
	g.Handle("/css/{file}", handler(cssHandler))

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", handler(newThreadHandler))
	p.Handle("/login/{id:[0-9]+}", handler(loginHandler))

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
