package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"carbocation.com/code/go.websocket-chat"
	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
	"github.com/carbocation/gotogether"
	//"github.com/garyburd/go-websocket/websocket"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// Master config, exported so it can be overrided
var Config *ConfigFile = Environment()

var (
	db     *sql.DB                                     //db maintains a pool of connections to our database of choice
	store  *sessions.FilesystemStore                   //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router *mux.Router               = mux.NewRouter() //Dynamic content is managed by handlers pointed at by the router
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

	wshub.Initialize(10*time.Second,
		60*time.Second,
		60*time.Second*9/10,
		4096,
		256)

	//Bundled static assets are handled by gotogether
	gotogether.Handle("/static/")

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", handler(indexHandler)).Name("index")
	g.Handle("/about", handler(aboutHandler)).Name("about")
	g.Handle("/forum/{id:[0-9]+}", handler(forumHandler)).Name("forum")
	g.Handle("/thread/{id:[0-9]+}", handler(threadHandler)).Name("thread")
	g.Handle("/thread", handler(newThreadHandler)).Name("newThread") //Form for creating new posts
	g.Handle("/login", handler(loginHandler)).Name("login")
	g.Handle("/logout", handler(logoutHandler)).Name("logout")
	g.Handle("/register", handler(registerHandler)).Name("register")
	g.HandleFunc(`/ws/thread/{id:[0-9]+}`, wsHandler)
	g.HandleFunc("/loaderio-3969952278183c9453e22d7f9ecfad1f/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "loaderio-3969952278183c9453e22d7f9ecfad1f")
	})

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", handler(postThreadHandler)).Name("postThread")
	p.Handle("/login", handler(postLoginHandler)).Name("postLogin")
	p.Handle("/register", handler(postRegisterHandler)).Name("postRegister")
	p.Handle("/vote", handler(postVoteHandler)).Name("postVote")

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
