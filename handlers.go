package main

import (
	"errors"
	"fmt"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/closuretable"
	"github.com/goods/httpbuf"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	/*
		//create the context
		ctx, err := NewContext(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer ctx.Close()
	*/

	//run the handler and grab the error, and report it
	buf := new(httpbuf.Buffer)
	//err = h(buf, req, ctx)
	err := h(buf, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//save the session
	/*
		if err = ctx.Session.Save(req, buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	*/

	//apply the buffered response to the writer
	buf.Apply(w)
}

func newThreadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	errors.New("Creating new threads is not yet implemented.")
	//http.Error(w, "Creating new threads is not yet implemented.\n", http.StatusInternalServerError)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	session, _ := store.Get(r, "user")
	defer session.Save(r, w)

	session.Values["id"] = mux.Vars(r)["id"]
	return
}

func defaultHandler(w http.ResponseWriter, r *http.Request) (err error) {
	remPartOfURL := r.URL.Path[len("/"):]
	fmt.Fprintf(w, "<html><head><link rel=\"stylesheet\" href=\"/css/main.css\"></head><body><h1>Welcome, %s</h1><a href='/hello/'>Say hello</a>", remPartOfURL)

	fmt.Fprint(w, "</body></html>")
	return
}

func cssHandler(w http.ResponseWriter, r *http.Request) (err error) {
	w.Header().Set("Content-Type", "text/css")

	file := mux.Vars(r)["file"]

	switch {
	case file == "main.css":
		fmt.Fprintf(w, "%s", mainCss())
	}

	return
}

func threadHandler(w http.ResponseWriter, r *http.Request) (err error) {
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
	
	return
}
