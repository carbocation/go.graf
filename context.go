package main

import (
	"encoding/gob"
	"net/http"

	"github.com/goods/httpbuf"
	"github.com/gorilla/context"
)

const (
	ThisUser = 0
)

func init() {
	//Tell gob about non-standard things we'll be serializing to disk
	gob.Register(new(User))
}

func OpenContext(req *http.Request) {
	session, _ := store.Get(req, "app")

	//Put the user into context
	if session.Values["user"] != nil {
		context.Set(req, ThisUser, session.Values["user"])
	} else {
		context.Set(req, ThisUser, new(User))
	}
}

func CloseContext(req *http.Request, buf *httpbuf.Buffer) (httpStatus int) {
	session, _ := store.Get(req, "app")

	err := session.Save(req, buf)

	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
