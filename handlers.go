package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
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
		G    *ConfigPublic
		User *user.User
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
	}
	//T("login.html").Execute(w, map[string]interface{}{})
	err = T("login.html").Execute(w, data)
	if err != nil {
		fmt.Printf("main.loginHandler: Template error: %s\n", err)
		return errors.New("Our template appears to be malformed so we cannot process your request.")
	}

	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) (err error) {
	DeleteContext(r, w)

	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return
}

func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	data := struct {
		G    *ConfigPublic
		User *user.User
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
	}

	err = T("index.html").Execute(w, data)
	if err != nil {
		fmt.Printf("main.indexHandler: Template error: %s\n", err)
		return errors.New("Our template appears to be malformed so we cannot process your request.")
	}

	return
}

func registerHandler(w http.ResponseWriter, r *http.Request) (err error) {
	data := struct {
		G        *ConfigPublic
		User     *user.User
		Messages []interface{}
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
		[]interface{}{},
	}

	//Don't let non-guests register again
	if !data.User.Guest() {
		http.Redirect(w, r, reverse("index"), http.StatusSeeOther)
	}

	session, _ := store.Get(r, "app")
	if flashes := session.Flashes(); len(flashes) > 0 {
		// Just print the flash values.
		data.Messages = flashes
	}

	err = T("register.html").Execute(w, data)
	if err != nil {
		fmt.Printf("main.registerHandler: Template error: %s\n", err)
		return errors.New("Our template appears to be malformed so we cannot process your request.")
	}
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
		fmt.Printf("main.threadHandler: %s\n", err)
		return errors.New("The requested thread's ancestry map could not be found.")
	}

	entries, err := forum.DescendantEntries(id)
	if err != nil {
		fmt.Printf("main.threadHandler: %s\n", err)
		return errors.New("The requested thread's neighbor entries could not be found.")
	}

	//Make sure this not a forum
	if entries[id].Forum {
		http.Redirect(w, r, reverse("forum", "id", id), http.StatusSeeOther)
		return
	}

	//Obligatory boxing step
	interfaceEntries := map[int64]interface{}{}
	for k, v := range entries {
		interfaceEntries[k] = v
	}

	tree, err := ct.TableToTree(interfaceEntries)
	if err != nil {
		fmt.Printf("main.threadHandler: Error converting closure table to tree: %s\n", err)
		return errors.New("The requested data structure could not be built.")
	}

	data := map[string]interface{}{
		"G":    Config.Public,
		"User": context.Get(r, ThisUser).(*user.User),
		"Tree": tree,
	}

	//execute the template
	err = T("thread.html").Execute(w, data)
	if err != nil {
		fmt.Printf("main.threadHandler: Template error: %s\n", err)
		return errors.New("Our template appears to be malformed so we cannot process your request.")
	}

	return
}

func forumHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//If the forum ID is not parseable as an integer, stop immediately
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return errors.New("The requested forum is invalid.")
	}

	// Pull down the closuretable from the root requested id
	ct, err := forum.DepthOneClosureTable(id)
	if err != nil {
		return errors.New("The requested forum's ancestry map could not be found.")
	}
	if ct.Size() < 1 {
		return errors.New("The requested forum's ancestry map had 0 entries.")
	}

	entries, err := forum.DepthOneDescendantEntries(id)
	if err != nil {
		return errors.New("The requested forum's neighbor entries could not be found.")
	}

	//Make sure this is a forum
	if entries[id].Forum != true {
		http.Redirect(w, r, reverse("thread", "id", id), http.StatusSeeOther)
		return
	}

	//Obligatory boxing step
	interfaceEntries := map[int64]interface{}{}
	for k, v := range entries {
		interfaceEntries[k] = v
	}

	tree, err := ct.TableToTree(interfaceEntries)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return errors.New("The requested data structure could not be built.")
	}

	data := map[string]interface{}{
		"G":    Config.Public,
		"User": context.Get(r, ThisUser).(*user.User),
		"Tree": tree,
	}

	//execute the template
	err = T("forum.html").Execute(w, data)
	if err != nil {
		fmt.Printf("main.forumHandler: Template error: %s\n", err)
		return errors.New("Our template appears to be malformed so we cannot process your request.")
	}

	return
}

func newThreadHandler(w http.ResponseWriter, r *http.Request) (err error) {
	fmt.Fprint(w, "New thread form will go here.")
	err = errors.New("The new thread form hasn't been created yet.")
	return
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	r.ParseForm()

	login := new(user.User)
	//Parse the form values into the Login object
	decoder.Decode(login, r.Form)

	u, err := login.Login()
	if err != nil {
		u = new(user.User)
	}

	context.Set(r, ThisUser, u)

	//Add the user's struct to the session
	session, _ := store.Get(r, "app")
	session.Values["user"] = u

	//Redirect to a GET address to prevent form resubmission
	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return
}

func postRegisterHandler(w http.ResponseWriter, r *http.Request) (err error) {
	r.ParseForm()

	//Don't let non-guests register again
	if !context.Get(r, ThisUser).(*user.User).Guest() {
		http.Redirect(w, r, reverse("index"), http.StatusSeeOther)
		return
	}

	//Make sure the repeat passwords match
	if r.FormValue("PlaintextPassword") != r.FormValue("PlaintextPassword2") {
		http.Redirect(w, r, reverse("register"), http.StatusSeeOther)
		return
	}

	//Locate the session
	session, _ := store.Get(r, "app")

	//Try to create the new user in the database
	u := new(user.User)
	decoder.Decode(u, r.Form)
	err = u.Register()
	if err != nil {
		//If our registration fails for any reason, set a flag and show the form again
		//http.Redirect(w, r, reverse("register"), http.StatusSeeOther)
		context.Set(r, ThisUser, u)

		//Tell the user why we failed
		session.AddFlash(fmt.Sprintf("%s", err))

		return registerHandler(w, r)
	}

	//They're a real user. Overwrite full object by populating from the DB
	u, err = user.FindOneUserById(u.Id)
	context.Set(r, ThisUser, u)

	session.Values["user"] = u

	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return
}

func postThreadHandler(w http.ResponseWriter, r *http.Request) error {
	var pid int64 //parent ID
	var err error
	var parent, entry *forum.Entry

	entry = new(forum.Entry)
	r.ParseForm()
	decoder.Decode(entry, r.Form)

	fmt.Printf("%+v", entry)

	//Don't let guests post (currently)
	//TODO(james) automatically create accounts for guests who try to post
	if context.Get(r, ThisUser).(*user.User).Guest() {
		http.Error(w, "NowayBro!", http.StatusUnauthorized)

		return errors.New("Unauthorized")
	}

	//TODO(james) stop relying on the existence of a user ID here
	user := context.Get(r, ThisUser).(*user.User)
	entry.AuthorId = user.Id

	//Make sure the parent_id is valid
	if pid, err = strconv.ParseInt(r.FormValue("parent_id"), 10, 64); err != nil {
		return err
	}

	//Make sure the parent post exists
	if parent, err = forum.OneEntry(pid); err != nil {
		return err
	}

	//If posting to a forum, must have a title and body or URL.
	if parent.Forum {
		//Min length
		if len(entry.Body) < 5 {
			if u, err := url.Parse(entry.Url); err != nil {
				return errors.New("Either free-text or a URL must be provided.")
			} else {
				//The URL is valid
				entry.Url = u.String()
			}
		}
		//entry.Body is therefore valid
	} else {
		//Thread reply
		entry.Url = ""
		if len(entry.Body) < 5 {
			return errors.New("Please craft a longer reply.")
		}
		//entry.Body is therefore valid
	}

	err = entry.Persist(parent.Id)
	if err != nil {
		return err
	}

	return nil
}
