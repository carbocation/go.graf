package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/carbocation/forum.git/forum"
	"github.com/goods/httpbuf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	session, _ := store.Get(r, "user")
	defer session.Save(r, w)

	session.Values["id"] = mux.Vars(r)["id"]
	return
}

type Demo struct {
	You string
}

func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//remPartOfURL := r.URL.Path[len("/"):]

	r.ParseForm()

	demo := new(Demo)
	decoder := schema.NewDecoder()
	decoder.Decode(demo, r.Form)

	//execute the template
	T("index.html").Execute(w, map[string]interface{}{
		"name": demo.You})

	return
}

func threadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//If the thread ID is not parseable as an integer, stop immediately
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return errors.New("The requested thread is invalid.")
	}

	// Pull down the closuretable from the root requested id
	ct, err := forum.ClosureTable(id)
	if err != nil {
		return errors.New("The requested thread could not be found.")
	}

	entries, err := forum.DescendantEntries(id)
	if err != nil {
		return errors.New("The requested thread could not be found.")
	}

	//Obligatory boxing step
	interfaceEntries := map[int64]interface{}{}
	for k, v := range entries {
		interfaceEntries[k] = v
	}

	tree, err := ct.TableToTree(interfaceEntries)
	if err != nil {
		return errors.New("The requested data structure could not be built.")
	}

	//execute the template
	T("thread.html").Execute(w, tree)

	return
}
