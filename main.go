package asksite

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"carbocation.com/code/go.websocket-chat"
	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
	"github.com/carbocation/gotogether"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// Master config, exported so it can be overrided
var Config *ConfigFile = Environment()

var (
	db     *sql.DB                                     //db maintains a pool of connections to our database of choice
	store  *sessions.FilesystemStore                   //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router *mux.Router               = mux.NewRouter() //Dynamic content is managed by Handlers pointed at by the router
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
	g.Handle("/", Handler(IndexHandler)).Name("index")
	g.Handle("/about", Handler(AboutHandler)).Name("about")
	g.Handle("/forum/{id:[0-9]+}", Handler(ForumHandler)).Name("forum")
	g.Handle("/thread/{id:[0-9]+}", Handler(ThreadHandler)).Name("thread")
	g.Handle("/thread", Handler(NewThreadHandler)).Name("newThread") //Form for creating new posts
	g.Handle("/login", Handler(LoginHandler)).Name("login")
	g.Handle("/logout", Handler(LogoutHandler)).Name("logout")
	g.Handle("/register", Handler(RegisterHandler)).Name("register")
	g.HandleFunc(`/ws/thread/{id:[0-9]+}`, ThreadWsHandler)

	//Create a subrouter for POST requests
	p := router.Methods("POST").Subrouter()
	p.Handle("/thread", Handler(PostThreadHandler)).Name("postThread")
	p.Handle("/login", Handler(PostLoginHandler)).Name("postLogin")
	p.Handle("/register", Handler(PostRegisterHandler)).Name("postRegister")
	p.Handle("/vote", Handler(PostVoteHandler)).Name("postVote")

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
