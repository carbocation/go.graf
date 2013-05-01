package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"runtime"

	"bitbucket.org/carbocation/nrsc"
	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// A config file is an object that nests various 
// public and other config structures
type ConfigFile struct {
	Public *ConfigPublic
}

type ConfigPublic struct {
	Site         string //Site name
	ContactEmail string //Webmaster email address
}

var (
	// Exported so it can be overrided
	Config *ConfigFile = &ConfigFile{
		Public: &ConfigPublic{
			Site:         "Ask Bitcoin",
			ContactEmail: "james@carbocation.com",
		},
	}

	db        *sql.DB     = initdb()                                                                                   //db maintains a pool of connections to our database of choice 
	appsecret             = `75Oop7MSN88WstKJSTyu9ALiO0Nbeckv/4/eDLDJcpXn0Ny1H9PdpzXDqApie77tZ04GFsdHehmzcMkAqh16Dg==` //64 bit random string generated with `openssl rand -base64 64`
	store                 = sessions.NewFilesystemStore("", []byte(appsecret))                                         //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router    *mux.Router = mux.NewRouter()                                                                            //Dynamic content is managed by handlers pointed at by the router 
)

// For exporting
func main() {
	//Only if we're running this package as the main package do we need to configure the maxprocs here
	runtime.GOMAXPROCS(runtime.NumCPU())

	//Call the main process
	Main()
}

func Main() {
	// Defer the close of the DB in the main function so we'll have a pool of connections maintained until the program exits
	defer db.Close()

	//Initialize the ancillary packages
	forum.Initialize(db)
	user.Initialize(db)

	//Bundled static assets are handled by nrsc
	nrsc.Handle("/static/")

	//Create a subrouter for GET requests
	g := router.Methods("GET").Subrouter()
	g.Handle("/", handler(indexHandler)).Name("index")
	g.Handle("/forum/{id:[0-9]+}", handler(forumHandler)).Name("forum")
	g.Handle("/thread/{id:[0-9]+}", handler(threadHandler)).Name("thread")
	g.Handle("/thread", handler(newThreadHandler)).Name("newThread") //Form for creating new posts
	g.Handle("/login", handler(loginHandler)).Name("login")
	g.Handle("/logout", handler(logoutHandler)).Name("logout")
	g.Handle("/register", handler(registerHandler)).Name("register")

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", handler(postThreadHandler)).Name("postThread")
	p.Handle("/login", handler(postLoginHandler)).Name("postLogin")
	p.Handle("/register", handler(postRegisterHandler)).Name("postRegister")

	//Notify the http package about our router
	http.Handle("/", router)

	//Launch the server
	if err := http.ListenAndServe("localhost:9999", nil); err != nil {
		panic(err)
	}
}

func initdb() *sql.DB {
	db, err := sql.Open("postgres", "dbname=projects user=askbitcoin password=xnkxglie port=5432 sslmode=disable")
	if err != nil {
		fmt.Println("Panic: " + err.Error())
		panic(err)
	}

	return db
}
