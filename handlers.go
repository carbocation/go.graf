package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/carbocation/forum.git/forum"
	"github.com/goods/httpbuf"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//Load session values into req	
	OpenContext(req)

	//For now, print the user's info to the console all the time
	fmt.Printf("User object: %+v\n", context.Get(req, ThisUser))

	//Run the handler and grab the error, and report it. We buffer the 
	// output so that handlers can modify session data at any point.
	buf := new(httpbuf.Buffer)
	if err := h(buf, req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Save any changed session values
	CloseContext(req, buf)

	//apply the buffered response to the writer
	buf.Apply(w)
}

func loginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//execute the template
	data := struct {
		G    GlobalValues
		User *User
	}{
		globals,
		context.Get(r, ThisUser).(*User),
	}
	//T("login.html").Execute(w, map[string]interface{}{})
	T("login.html").Execute(w, data)
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) (err error) {
	DeleteContext(r, w)

	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	data := struct {
		G    GlobalValues
		User *User
	}{
		globals,
		context.Get(r, ThisUser).(*User),
	}

	T("index.html").Execute(w, data)

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

	data := map[string]interface{}{
		"G":    globals,
		"User": context.Get(r, ThisUser).(*User),
		"Tree": tree,
	}

	//execute the template
	T("thread.html").Execute(w, data)

	return
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	r.ParseForm()

	login := new(Login)
	//Parse the form values into the Login object
	decoder.Decode(login, r.Form)

	user, err := login.Login()
	if err != nil {

		//They're a guest user
		context.Set(r, ThisUser, &User{})
	} else {
		//They're a real user
		context.Set(r, ThisUser, user)
	}

	//Add the user's struct to the session
	session, _ := store.Get(r, "app")
	session.Values["user"] = user

	//Redirect to a GET address to prevent form resubmission
	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return
}

func postThreadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	errors.New("Creating new threads is not yet implemented.")
	return
}
