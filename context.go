package graf

import (
	"encoding/gob"
	"net/http"

	"github.com/carbocation/go.user"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
)

const (
	ThisUser = 0
)

func init() {
	//Tell gob about non-standard things we'll be serializing to disk
	gob.Register(new(user.User))
}

func OpenContext(req *http.Request) {
	session, _ := store.Get(req, "app")

	//Put the user into context
	if session.Values["user"] != nil {
		context.Set(req, ThisUser, session.Values["user"])
	} else {
		context.Set(req, ThisUser, new(user.User))
	}
}

//Anything that satisfies the http.ResponseWriter interface is sufficient
func CloseContext(req *http.Request, buf http.ResponseWriter) (httpStatus int) {
	session, _ := store.Get(req, "app")

	err := session.Save(req, buf)

	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func DeleteContext(req *http.Request, w http.ResponseWriter) {
	//Delete all user-set variables
	context.Purge(0)

	//Destroy the session
	session, _ := store.Get(req, "app")
	session.Options = &sessions.Options{MaxAge: -1}
	session.Save(req, w)
}
