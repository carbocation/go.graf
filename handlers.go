package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	//"html/template"

	//"github.com/carbocation/gotogether"
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

//For now, the index is actually just a hardlink to the forum with ID #1
func indexHandler(w http.ResponseWriter, r *http.Request) (err error) {
	mux.Vars(r)["id"] = "1"

	return forumHandler(w, r)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) (err error) {
	data := struct {
		G    *ConfigPublic
		User *user.User
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
	}

	err = T("about.html").Execute(w, data)
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
	
	u := context.Get(r, ThisUser).(*user.User)

	tree, err := forum.DescendantEntries(id, u)
	if err != nil {
		fmt.Printf("main.threadHandler: %s\n", err)
		return errors.New("The requested thread's neighbor entries could not be found.")
	}

	//Make sure this not a forum
	if tree.Forum {
		http.Redirect(w, r, reverse("forum", "id", id), http.StatusSeeOther)
		return
	}

	data := struct {
		G    *ConfigPublic
		User *user.User
		Tree *forum.Entry
	}{
		G:    Config.Public,
		User: u,
		Tree: tree,
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
	
	u := context.Get(r, ThisUser).(*user.User)

	tree, err := forum.DescendantEntries(id, u)
	if err != nil {
		return errors.New("The requested forum's neighbor entries could not be found.")
	}

	//Make sure this is a forum
	if tree.Forum != true {
		http.Redirect(w, r, reverse("thread", "id", id), http.StatusSeeOther)
		return
	}

	data := struct {
		G    *ConfigPublic
		User *user.User
		Tree *forum.Entry
	}{
		G:    Config.Public,
		User: u,
		Tree: tree,
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
	
	//Don't let guests post (currently)
	//TODO(james) automatically create accounts for guests who try to post
	u := context.Get(r, ThisUser).(*user.User)
	if u.Guest() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return errors.New("Only logged-in users can post.")
	}
	
	//Make sure the parent_id is valid
	if pid, err = strconv.ParseInt(r.FormValue("parent_id"), 10, 64); err != nil {
		return err
	}

	//Make sure the parent post exists
	if parent, err = forum.OneEntry(pid); err != nil {
		return err
	}

	entry = new(forum.Entry)
	r.ParseForm()
	decoder.Decode(entry, r.Form)
	
	entry.AuthorId = u.Id
	entry.Body = strings.TrimSpace(entry.Body)
	
	URL, err := url.Parse(strings.TrimSpace(r.FormValue("Url")))
	if err != nil {
		return err
	}

	if URL.String() == "" && entry.Body == "" {
		//Lack of URL and Body fails in all contexts
		http.Error(w, "Unauthorized", http.StatusExpectationFailed)
		return errors.New("Both URL and Body were empty; please fill out either one or the other.")
	}
	
	//When creating new posts, we set their parent to their true parent (we don't use 
	// LCRS at that stage), so checking for Parent().Forum is sufficient.
	if parent.Forum && entry.Title == "" {
		//Unacceptable to have an empty title if this is new entry within a forum
		return errors.New("The Title must not be empty or consist solely of whitespace.")
	} 
	
	if parent.Forum && ValidUrl(URL) {
		//We only care about whether the parent is a forum if the user submits a valid URL
		//As promised, we replace the Body with the URL if one is given 
		entry.Body = URL.String()
		entry.Url = true
	} else if len(entry.Body) < 5 {
		//In all other cases, if the body is not valid, they need to write more.
		return errors.New("Please craft a longer message.")
	}
	
	fmt.Printf("Entry just before persistence: %v", entry)

	err = entry.Persist(parent.Id)
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(entry)
	integer, err := w.Write(jsondata)
	if err != nil {
		return err
	}

	fmt.Printf("Integer from posting the new entry was %d\n", integer)
	//We can set the content type after sending the jsondata because
	// we're actually using buffered output
	w.Header().Set("Content-type", "application/json")

	return nil
}

func postVoteHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	//Make sure the target entry is valid
	entryId, err := strconv.ParseInt(r.FormValue("entryId"), 10, 64)
	if err != nil {
		return err
	}

	entry, err := forum.OneEntry(entryId)
	if err != nil {
		return err
	}

	//Don't let guests post (currently)
	//TODO(james) automatically create accounts for guests who try to post
	if context.Get(r, ThisUser).(*user.User).Guest() {
		http.Error(w, "NowayBro!", http.StatusUnauthorized)

		return errors.New("Unauthorized")
	}

	//TODO(james) stop relying on the existence of a user ID here
	user := context.Get(r, ThisUser).(*user.User)

	vote := &forum.Vote{EntryId: entry.Id, UserId: user.Id}

	if r.FormValue("vote") == "upvote" {
		vote.Upvote, vote.Downvote = true, false
	} else if r.FormValue("vote") == "downvote" {
		vote.Upvote, vote.Downvote = false, true
	} else {
		vote.Upvote, vote.Downvote = false, false
	}

	err = vote.Persist()
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(vote)

	integer, err := w.Write(jsondata)
	if err != nil {
		return err
	}

	fmt.Printf("Integer from posting the new entry was %d\n", integer)
	//We can set the content type after sending the jsondata because
	// we're actually using buffered output
	w.Header().Set("Content-type", "application/json")

	return nil
}
