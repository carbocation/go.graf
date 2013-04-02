package main

import (
	"fmt"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"github.com/carbocation/util.git/datatypes/closuretable"
	"net/http"
	"strconv"
	//"html/template"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var db *sql.DB

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func main() {
	// Initialize the DB in the main function so we'll have a pool of connections maintained
	db = initdb()
	defer db.Close()

	//Initialize our router
	r := mux.NewRouter()

	//Create a subrouter for GET requests
	g := r.Methods("GET").Subrouter()
	g.HandleFunc("/", defaultHandler)
	g.HandleFunc("/thread/{id:[0-9]+}", threadHandler)
	g.HandleFunc("/css/{file}", cssHandler)

	//Create a subrouter for POST requests
	p := r.Methods("POST").Subrouter()
	p.HandleFunc("/thread", newThreadHandler)
	p.HandleFunc("/login/{id:[0-9]+}", loginHandler)

	//Notify the http package about our router
	http.Handle("/", r)

	//Launch the server
	http.ListenAndServe("localhost:9999", nil)
}

func initdb() *sql.DB {
	db, err := sql.Open("postgres", "dbname=forumtest sslmode=disable")
	if err != nil {
		panic(err)
	}

	return db
}

func newThreadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creating new threads is not yet implemented.\n")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user")
	defer session.Save(r, w)

	session.Values["id"] = mux.Vars(r)["id"]
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/"):]
	fmt.Fprintf(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body><h1>Welcome, %s</h1><a href='/hello/'>Say hello</a>", remPartOfURL)

	fmt.Fprint(w, "</body></html>")
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	file := mux.Vars(r)["file"]

	switch {
	case file == "main.css":
		fmt.Fprintf(w, "%s", mainCss())
	}

}

func mainCss() string {
	return `
div .comment {
	padding-left: 100px;
}
`
}

func threadHandler(w http.ResponseWriter, r *http.Request) {
	unsafeId := r.URL.Path[len("/thread/"):]

	//If the thread ID is not parseable as an integer, stop immediately
	id, err := strconv.ParseInt(unsafeId, 10, 64)
	if err != nil {
		return
	}

	// Generate a closuretable from the root requested id
	ct := closuretable.New(id)
	// Pull down the remaining elements in the closure table that are descendants of this node
	q := `select * 
from entry_closures
where descendant in (
select descendant
from entry_closures
where ancestor=$1
)
and ancestor in (
select descendant
from entry_closures
where ancestor=$1
)
and depth = 1`
	stmt, err := db.Prepare(q)
	if err != nil {
		//fmt.Printf("Statement Preparation Error: %s", err)
		return
	}

	rows, err := stmt.Query(unsafeId)
	if err != nil {
		//fmt.Printf("Query Error: %v", err)
		return
	}

	//Populate the closuretable
	for rows.Next() {
		var ancestor, descendant int64
		var depth int
		err = rows.Scan(&ancestor, &descendant, &depth)
		if err != nil {
			//fmt.Printf("Rowscan error: %s", err)
			return
		}

		err = ct.AddChild(closuretable.Child{Parent: ancestor, Child: descendant})

		//err = ct.AddRelationship(closuretable.Relationship{Ancestor: ancestor, Descendant: descendant, Depth: depth})
		if err != nil {
			//fmt.Fprintf(w, "Error: %s", err)
			return
		}
	}

	id, entries, err := forum.RetrieveDescendantEntries(unsafeId, db)
	if err != nil {
		//fmt.Fprintf(w, "Error: %s", err)
		return
	}

	//Obligatory boxing step
	interfaceEntries := map[int64]interface{}{}
	for k, v := range entries {
		interfaceEntries[k] = v
	}

	tree, err := ct.TableToTree(interfaceEntries)
	if err != nil {
		//fmt.Printf("TableToTree error: %s", err)
		return
	}

	//Spew the posts' HTML over a channel
	htm := make(chan string)
	go func() {
		PrintNestedComments(tree, htm)
		close(htm)
	}()

	fmt.Fprint(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body>")
	for h := range htm {
		fmt.Fprint(w, h)
	}
	fmt.Fprint(w, "</body></html>")
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
