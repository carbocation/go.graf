package graf

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/carbocation/go.websocket-chat"
	"github.com/carbocation/go.forum"
	"github.com/carbocation/go.user"
	"github.com/garyburd/go-websocket/websocket"
	"github.com/goods/httpbuf"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Handler func(http.ResponseWriter, *http.Request) error

/*
Derived from https://github.com/gorilla/handlers/blob/master/handlers.go
The main purpose here is to be able to pass around http.Status codes, which
the net/http package does not natively enable.
*/
type ResponseLogger struct {
	httpbuf.Buffer
	size   int //Size of buffer
	status int //HTTP Status Code
}

//Get the current status code
func (l ResponseLogger) Status() int {
	return l.status
}

//Set a status code. ints are accepted, but
// http.StatusXXX is easier to read. Your choice.
func (l *ResponseLogger) WriteHeader(s int) {
	l.status = s
	l.Buffer.WriteHeader(s)
}

//Call Status before you call this
func (l *ResponseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.Buffer.Write(b)
	l.size += size
	return size, err
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//Load any session values into req
	OpenContext(req)

	//Run the handler and grab the error, and report it. We buffer the
	//output so that handlers can modify session/header data at any point.
	buf := new(ResponseLogger)
	err := h(buf, req)

	//Log the request
	size, status, username := buf.Len(), buf.Status(), "-"
	if uifc := context.Get(req, ThisUser); uifc != nil {
		//Provide the real username, if any
		u := uifc.(*user.User)
		if !u.Guest() {
			username = u.Handle
		}
	}

	LogWriter.Print(fmt.Sprintf(`%s - "%s" [%s] "%s %s %s" %d %d "%s" "%s"`,
		strings.Split(getIpAddress(req), ":")[0],
		username,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		req.Method,
		req.RequestURI,
		req.Proto,
		status,
		size,
		req.Referer(),
		req.UserAgent(),
	))

	if err != nil {
		//Errors for which there should be HTML output should be invoked by
		//calling and returning ErrorHTML() in the handler, rather than just
		//by returning an error. This should be used to catch any non-formatted
		//error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)

		ErrorLogWriter.Print(fmt.Sprintf(`%s - [%s] "%s %s %s" "%s" "%s" "%s"`,
			strings.Split(getIpAddress(req), ":")[0],
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Method,
			req.RequestURI,
			req.Proto,
			req.Referer(),
			req.UserAgent(),
			err.Error(),
		))

		return
	}

	//Save any changed session values
	CloseContext(req, buf)

	//Apply the buffered response to the actual writer
	buf.Apply(w)
}

//From https://groups.google.com/forum/?fromgroups#!topic/golang-nuts/lomWKs0kOfE
func getIpAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIp := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIp == "" && hdrForwardedFor == "" {
		return r.RemoteAddr
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	return hdrRealIp
}

//Produce an HTML error page based on a title and a message, and return a desired error code.
//It is encouraged but not mandatory to use http.StatusXXX codes instead of raw integers for errorCode
func ErrorHTML(w http.ResponseWriter, r *http.Request, errorTitle, errorMessage string, errorCode int, internalError error) error {
	w.WriteHeader(errorCode)

	//Log it
	ErrorLogWriter.Print(fmt.Sprintf(`%s - [%s] "%s %s %s" "%s" "%s" "%s"`,
		strings.Split(getIpAddress(r), ":")[0],
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method,
		r.RequestURI,
		r.Proto,
		r.Referer(),
		r.UserAgent(),
		internalError.Error(),
	))

	data := struct {
		G          *ConfigPublic
		User       *user.User
		ShortError string
		LongError  string
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
		errorTitle,
		errorMessage,
	}

	err := T("error.html").Execute(w, data)
	if err != nil {
		return errors.New("First, there was an error. Then, there was an error in our tool that handles errors. Therefore, we doubly cannot complete your request.")
	}

	return err
}

func LoginHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//execute the template
	data := struct {
		G        *ConfigPublic
		User     *user.User
		Messages []interface{}
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
		[]interface{}{},
	}

	session, _ := store.Get(r, "app")
	if flashes := session.Flashes(); len(flashes) > 0 {
		// Just print the flash values.
		data.Messages = flashes
	}

	//T("login.html").Execute(w, map[string]interface{}{})
	err = T("login.html").Execute(w, data)
	if err != nil {
		return ErrorHTML(w, r, "Bad template", "We cannot complete your request because we are having problems with our templating engine.", http.StatusInternalServerError, err)
	}

	return
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	DeleteContext(r, w)

	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return nil
}

//For now, the index is actually just a hardlink to the forum with ID #1
func IndexHandler(w http.ResponseWriter, r *http.Request) error {
	mux.Vars(r)["id"] = Config.App.RootForumID

	return ForumHandler(w, r)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) error {
	data := struct {
		G    *ConfigPublic
		User *user.User
	}{
		Config.Public,
		context.Get(r, ThisUser).(*user.User),
	}

	err := T("about.html").Execute(w, data)
	if err != nil {
		return ErrorHTML(w, r, "Bad template", "We cannot complete your request because we are having problems with our templating engine.", http.StatusInternalServerError, err)
	}

	return err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) (err error) {
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

		return nil
	}

	session, _ := store.Get(r, "app")
	if flashes := session.Flashes(); len(flashes) > 0 {
		// Just add the flash values.
		data.Messages = flashes
	}

	err = T("register.html").Execute(w, data)
	if err != nil {
		return ErrorHTML(w, r, "Bad template", "We cannot complete your request because we are having problems with our templating engine.", http.StatusInternalServerError, err)
	}
	return
}

func ThreadHandler(w http.ResponseWriter, r *http.Request) error {
	//If the thread ID is not parseable as an integer, stop immediately
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return ErrorHTML(w, r, "Invalid thread", "The requested thread is invalid.", http.StatusBadRequest, err)
	}

	u := context.Get(r, ThisUser).(*user.User)

	tree, err := forum.DescendantEntries(id, u)
	if err != nil {
		return ErrorHTML(w, r, "Thread error", "The request could not be completed due to an internal server error.", http.StatusInternalServerError, err)
	}

	if tree == nil {
		return ErrorHTML(w, r, "Thread not found", "The thread that you requested could not be found.", http.StatusNotFound, errors.New("tree was nil"))
	}

	//Make sure this not a forum
	if tree.Forum {
		http.Redirect(w, r, reverse("forum", "id", id), http.StatusSeeOther)
		return nil
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
		return ErrorHTML(w, r, "Bad template", "We cannot complete your request because we are having problems with our templating engine.", http.StatusInternalServerError, err)
	}

	return nil
}

func ForumHandler(w http.ResponseWriter, r *http.Request) (err error) {
	//If the forum ID is not parseable as an integer, stop immediately
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		return ErrorHTML(w, r, "Invalid forum", "The requested forum is invalid.", http.StatusBadRequest, err)
	}

	u := context.Get(r, ThisUser).(*user.User)

	tree, err := forum.DescendantEntries(id, u)
	if err != nil {
		return ErrorHTML(w, r, "Forum error", "The request could not be completed due to an internal server error.", http.StatusInternalServerError, err)
	}

	if tree == nil {
		return ErrorHTML(w, r, "Forum not found", "The forum that you requested could not be found.", http.StatusNotFound, errors.New("tree was nil"))
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
		return ErrorHTML(w, r, "Template error", "We had a templating malfunction and could not serve your request.", http.StatusInternalServerError, err)
	}

	return
}

//Handles the specific URL /thread without any ID affixed to it
func ThreadRootHandler(w http.ResponseWriter, r *http.Request) (err error) {
	
	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)
	
	return
}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	session, _ := store.Get(r, "app")

	login := new(user.User)
	//Parse the form values into the Login object
	decoder.Decode(login, r.Form)

	u, err := login.Login()
	if err != nil {
		//Add a flash error message stating the error
		session.AddFlash(fmt.Sprintf("%s", err))

		//Write the http error code
		w.WriteHeader(http.StatusBadRequest)

		//Send them back to the login form
		return LoginHandler(w, r)
	}

	//Successful login

	context.Set(r, ThisUser, u)

	//Add the user's struct to the session
	session.Values["user"] = u

	//Redirect to a GET address to prevent form resubmission
	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return nil
}

func PostRegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var err error
	r.ParseForm()

	//Don't let non-guests register again
	if !context.Get(r, ThisUser).(*user.User).Guest() {
		http.Redirect(w, r, reverse("index"), http.StatusSeeOther)
		return err
	}

	//Make sure the repeat passwords match
	if r.FormValue("PlaintextPassword") != r.FormValue("PlaintextPassword2") {
		http.Redirect(w, r, reverse("register"), http.StatusSeeOther)
		return err
	}

	//Locate the session
	session, _ := store.Get(r, "app")

	//Try to create the new user in the database
	u := new(user.User)
	decoder.Decode(u, r.Form)
	err = u.Register()
	if err != nil {
		//If our registration fails for any reason, set a flag and show the form again
		context.Set(r, ThisUser, u)

		//Tell the user why we failed
		session.AddFlash(fmt.Sprintf("%s", err))

		//Tell the browser that that input was no good
		w.WriteHeader(http.StatusBadRequest)

		return RegisterHandler(w, r)
	}

	//They're a real user. Overwrite full object by populating from the DB
	u, err = user.FindOneUserById(u.Id)
	context.Set(r, ThisUser, u)

	session.Values["user"] = u

	http.Redirect(w, r, reverse("index"), http.StatusSeeOther)

	return err
}

func PostThreadHandler(w http.ResponseWriter, r *http.Request) error {
	var pid int64 //parent ID
	var err error
	var parent, entry *forum.Entry

	//Don't let guests post (currently)
	//TODO(james) automatically create accounts for guests who try to post
	u := context.Get(r, ThisUser).(*user.User)
	if u.Guest() {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Only logged-in users can post.")
	}

	//Make sure the parent_id is valid
	if pid, err = strconv.ParseInt(r.FormValue("parent_id"), 10, 64); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	//Make sure the parent post exists
	if parent, err = forum.OneEntry(pid); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Both URL and Body were empty; please fill out either one or the other.")
	}

	//When creating new posts, we set their parent to their true parent (we don't use
	// LCRS at that stage), so checking for Parent().Forum is sufficient.
	if parent.Forum && entry.Title == "" {
		//Unacceptable to have an empty title if this is new entry within a forum
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("The Title must not be empty or consist solely of whitespace.")
	}

	if parent.Forum && ValidUrl(URL) {
		//We only care about whether the parent is a forum if the user submits a valid URL
		//As promised, we replace the Body with the URL if one is given
		entry.Body = URL.String()
		entry.Url = true
	} else if len(entry.Body) < 5 {
		//In all other cases, if the body is not valid, they need to write more.
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Please craft a longer message.")
	}

	//TODO(james): clean up this parent.Id and ParentId junk
	err = entry.Persist(parent.Id)
	if err != nil {
		return err
	}
	
	//When you post, you give yourself an upvote
	vote := &forum.Vote{EntryId: entry.Id, UserId: u.Id, Upvote: true, Downvote: false}
	err = vote.Persist()
	if err != nil {
		return err
	}
	entry.UserVote = vote
	entry.Upvotes = 1
	//End upvote

	//Intercept here to send the data down to websocket listeners
	go func() {
		//OVERWRITING ENTRY WITH THE NEWLY-MINTED VERSION FROM THE DATABASE
		entry, err = forum.OneEntry(entry.Id)
		if err != nil {
			return
		}
		entry.ParentId = parent.Id

		//Identify all ancestor posts so anyone viewing them directly will get 
		//notified of this new post via Websockets
		e, _ := forum.AncestorEntries(entry.Id, u)
		ids := []string{}
		for e != nil {
			ids = append(ids, strconv.Itoa(int(e.Id)))
			fmt.Println(e)
			e = e.Child()
		}

		packet, _ := wshub.Packetize("thread_post", *entry)
		wshub.Multicast(packet, ids)
	}()

	//TODO(james): Delete the rest?

	jsondata, err := json.Marshal(entry)
	_, err = w.Write(jsondata)
	if err != nil {
		return err
	}

	//We can set the content type after sending the jsondata because
	// we're actually using buffered output
	w.Header().Set("Content-type", "application/json")

	return nil
}

func PostVoteHandler(w http.ResponseWriter, r *http.Request) error {
	r.ParseForm()

	//Make sure the target entry is valid
	entryId, err := strconv.ParseInt(r.FormValue("entryId"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	entry, err := forum.OneEntry(entryId)
	if err != nil {
		return err
	}

	//Don't let guests post (currently)
	//TODO(james) automatically create accounts for guests who try to post
	if context.Get(r, ThisUser).(*user.User).Guest() {
		w.WriteHeader(http.StatusBadRequest)

		return errors.New("Only logged-in users can cast votes.")
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

	_, err = w.Write(jsondata)
	if err != nil {
		return err
	}

	//We can set the content type after sending the jsondata because
	// we're actually using buffered output
	w.Header().Set("Content-type", "application/json")

	return nil
}

func ThreadWsHandler(w http.ResponseWriter, req *http.Request) {
	ws, err := websocket.Upgrade(w, req.Header, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	//When we try to handle this, see if the hub exists.
	id := mux.Vars(req)["id"]
	wshub.Launch(ws, id)
}
