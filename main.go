package main

import (
	"fmt"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"github.com/carbocation/util.git/datatypes/closuretable"
	"net/http"
	"strconv"
	"time"
	//"html/template"
	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/thread/", threadHandler)
	http.HandleFunc("/hello/", commentHandler)
	http.HandleFunc("/css/", cssHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe("localhost:9999", nil)
}

func initdb() *sql.DB {
	db, err := sql.Open("postgres", "dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	
	return db
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/"):]
	fmt.Fprintf(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body><h1>Welcome, %s</h1><a href='/hello/'>Say hello</a>", remPartOfURL)

	fmt.Fprint(w, "</body></html>")
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	docname := r.URL.Path[len("/css/"):]

	switch {
	case docname == "main.css":
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

func commentHandler(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/hello/"):] //get everything after the /hello/ part of the URL
	//w.Header().Set("Content-Type", "text/html")

	fmt.Fprint(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body>")
	fmt.Fprintf(w, "Hello %s!", remPartOfURL)

	PrintNestedComments(w, ClosureTree())

	fmt.Fprint(w, "</body></html>")
}

func threadHandler(w http.ResponseWriter, r *http.Request) {
	db := initdb()
	defer db.Close()
	
	unsafeId := r.URL.Path[len("/thread/"):] //get everything after the /hello/ part of the URL
	
	// Prepare a statement
	stmt, err := db.Prepare("SELECT * FROM golang.yourtable WHERE rank < $1")
	if err != nil {
		fmt.Printf("Statement Preparation Error: %s", err)
	}

	// Query from that prepared statement
	rows, err := stmt.Query(5)
	if err != nil {
		fmt.Printf("Query Error: %v", err)
	}

	// Show the *Rows ptr
	fmt.Fprintf(w, "Row pointer: %#v \n", rows)

	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("Column error: %s", err)
	}

	fmt.Fprintf(w, "Columns: %s \n", cols)

	// Iterate over the rows
	for rows.Next() {
		var rank int
		var username, password string
		err = rows.Scan(&rank, &username, &password)
		fmt.Fprintf(w, "Record: %#i, %s, %s \n", rank, username, password)
	}

	//If the thread ID is not parseable as an integer, stop them
	id, err := strconv.Atoi(unsafeId)
	if err != nil {
		// Default to the root if they gave us a non-integer value
		id = 0
	}

	tree := ClosureTree()

	fmt.Fprint(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body>")
	//Print the ID they gave us
	fmt.Fprintf(w, "Hello %i, here is a tree %+v", id, tree)
	fmt.Fprint(w, "</body></html>")
}

func PrintNestedComments(w http.ResponseWriter, el *binarytree.Tree) {
	if el == nil {
		return
	}

	fmt.Fprint(w, "<div class=\"comment\">")
	//Self
	e := el.Value.(forum.Entry)
	fmt.Fprintf(w, "Title: %s", e.Title)

	//Children are nested
	PrintNestedComments(w, el.Left())
	fmt.Fprint(w, "</div>")

	//Siblings are parallel
	PrintNestedComments(w, el.Right())
}

func ClosureTree() *binarytree.Tree {
	//Make some entries
	entries := map[int64]forum.Entry{
		0: forum.Entry{Id: 100, Title: "Title 100", Body: "Body 100", Created: time.Now(), AuthorId: 0},
		1: forum.Entry{Id: 101, Title: "Title 101", Body: "Body 101", Created: time.Now(), AuthorId: 1},
		2: forum.Entry{Id: 102, Title: "Title 102", Body: "Body 102", Created: time.Now(), AuthorId: 2},
		3: forum.Entry{Id: 103, Title: "Title 103", Body: "Body 103", Created: time.Now(), AuthorId: 3},
	}

	ct := closureTable()

	// Obligatory boxing step
	// Convert to interface type so the generic TableToTree method can be called on these entries
	boxedEntries := map[int64]interface{}{}
	for k, v := range entries {
		boxedEntries[k] = v
	}

	tree := ct.TableToTree(boxedEntries)

	return tree
}

func closureTable() *closuretable.ClosureTable {
	//Make a hierarchy for these entries
	ct := closuretable.New(0)
	ct.AddChild(closuretable.Child{Parent: 0, Child: 1})
	ct.AddChild(closuretable.Child{Parent: 0, Child: 2})
	ct.AddChild(closuretable.Child{Parent: 1, Child: 3})

	return ct
}
