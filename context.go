package main

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type Context struct {
	Session *sessions.Session
	User    *User
}

func (c *Context) Close() {

}

func NewContext(req *http.Request) (ctx *Context, err error) {
	sess, err := store.Get(req, "app")
	ctx = &Context{
		Session: sess,
	}
	if err != nil {
		return
	}

	//try to fill in the user from the session
	/*
		if uid, ok := sess.Values["user"].(bson.ObjectId); ok {
			err = ctx.C("users").Find(bson.M{"_id": uid}).One(&ctx.User)
		}
	*/

	return
}
