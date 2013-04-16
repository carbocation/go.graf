package main
/*
import (
	"net/http"
)

func NewContext(req *http.Request) (ctx *Context, err error) {
	sess, err := store.Get(req, "gostbook")
	ctx = &Context{
		Session:  sess,
	}
	if err != nil {
		return
	}

	//try to fill in the user from the session
	if uid, ok := sess.Values["user"].(bson.ObjectId); ok {
		err = ctx.C("users").Find(bson.M{"_id": uid}).One(&ctx.User)
	}

	return
}
*/