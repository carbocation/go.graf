/*
This is an example application that puts together all of the pieces
to make use of the AskSite toolkit for building threaded forums in Golang.

Copyright 2013 James Pirruccello <james@carbocation.com>
Reproduction and use are governed by the terms of the LICENSE file in this folder
*/
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"carbocation.com/code/go.gtfo"
	"carbocation.com/code/go.websocket-chat"
	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
	"github.com/carbocation/gotogether"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var (
	Config  *asksite.ConfigFile       = Environment()   // Master config, exported so it can be overrided
	db      *sql.DB                                     //db maintains a pool of connections to our database of choice
	store   *sessions.FilesystemStore                   //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router  *mux.Router               = mux.NewRouter() //Dynamic content is managed by asksite.Handlers pointed at by the router
	decoder *schema.Decoder           = schema.NewDecoder()
)

func main() {
	//
	//After user has had opportunity to change config:
	//
	//1 init the db
	db = initdb()
	// Defer the close of the DB in the main function so we'll have a pool of connections maintained until the program exits
	defer db.Close()

	//2 setup our session store
	store = sessions.NewFilesystemStore("", []byte(Config.App.Secret))

	//Initialize the ancillary packages
	forum.Initialize(db)
	user.Initialize(db)
	asksite.Initialize(Config, db, store, router, decoder)

	wshub.Initialize(10*time.Second,
		60*time.Second,
		60*time.Second*9/10,
		4096,
		256)

	//Bundled static assets are handled by gotogether
	gotogether.Handle("/static/")

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", asksite.Handler(asksite.IndexHandler)).Name("index")
	g.Handle("/about", asksite.Handler(asksite.AboutHandler)).Name("about")
	g.Handle("/forum/{id:[0-9]+}", asksite.Handler(asksite.ForumHandler)).Name("forum")
	g.Handle("/thread/{id:[0-9]+}", asksite.Handler(asksite.ThreadHandler)).Name("thread")
	g.Handle("/thread", asksite.Handler(asksite.NewThreadHandler)).Name("newThread") //Form for creating new posts
	g.Handle("/login", asksite.Handler(asksite.LoginHandler)).Name("login")
	g.Handle("/logout", asksite.Handler(asksite.LogoutHandler)).Name("logout")
	g.Handle("/register", asksite.Handler(asksite.RegisterHandler)).Name("register")
	g.HandleFunc(`/ws/thread/{id:[0-9]+}`, asksite.ThreadWsHandler)
	g.HandleFunc("/loaderio-38b140f4cb51d3ffca9d71c7529a336d/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "loaderio-38b140f4cb51d3ffca9d71c7529a336d")
	})

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", asksite.Handler(asksite.PostThreadHandler)).Name("postThread")
	p.Handle("/login", asksite.Handler(asksite.PostLoginHandler)).Name("postLogin")
	p.Handle("/register", asksite.Handler(asksite.PostRegisterHandler)).Name("postRegister")
	p.Handle("/vote", asksite.Handler(asksite.PostVoteHandler)).Name("postVote")

	//Notify the http package about our router
	http.Handle("/", router)

	//Launch the server
	if err := http.ListenAndServe(fmt.Sprintf("localhost:%s", Config.App.Port), nil); err != nil {
		panic(err)
	}
}

func initdb() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s port=%s sslmode=disable",
		Config.DB.DBName,
		Config.DB.User,
		Config.DB.Password,
		Config.DB.Port))
	db.SetMaxIdleConns(Config.DB.PoolSize)
	if err != nil {
		fmt.Println("Panic: " + err.Error())
		panic(err)
	}

	return db
}
