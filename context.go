package main

import (
	//"github.com/gorilla/context"
)

const (
	ThisUser = 0
)

/*
import (
	"fmt"
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
	if err != nil {
		return
	}
	ctx = &Context{
		Session: sess,
	}

	//testing
	fmt.Printf("%+v", sess.Values)
	//testing

	var uid int64 = 0
	if ctx.Session.Values["id"] != nil {
		uid = ctx.Session.Values["id"].(int64)
	}

	//try to fill in the user from the session
	ctx.User, err = FindOneUserById(uid)
	/*
		if uid, ok := sess.Values["user"].(bson.ObjectId); ok {
			err = ctx.C("users").Find(bson.M{"_id": uid}).One(&ctx.User)
		}


	return
}
*/
