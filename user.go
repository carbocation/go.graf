package main

import (
	"errors"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

type Password string

type User struct {
	Id       int64     "The user's auto-incremented ID"
	Handle   string    "The user's name"
	Password Password  "The byteslice of the user's bcrypted password"
	Created  time.Time "The creation timestamp of the user's account"
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password Password) (err error) {
	if len(password) < 1 {
		err = errors.New("Error: no password was chosen.")
		return
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = Password(hpass)
	
	return
}

//Login validates and returns a user object if they exist in the database.
/*
//Commented out because we have no *Context so far.
func Login(ctx *Context, username, password string) (u *User, err error) {
	err = ctx.C("users").Find(bson.M{"username": username}).One(&u)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	if err != nil {
		u = nil
	}
	return
}
*/
