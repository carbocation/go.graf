package main

import (
	"errors"
	//"fmt"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
)

type Login struct {
	Handle            string "The username"
	PlaintextPassword string "The password"
}

type User struct {
	Id       int64     "The user's auto-incremented ID"
	Handle   string    "The user's name"
	Password string    "The user's bcrypted password"
	Created  time.Time "The creation timestamp of the user's account"
}

//SetPassword takes a plaintext password and hashes it with bcrypt and sets the
//password field to the hash.
func (u *User) SetPassword(password string) (err error) {
	if len(password) < 1 {
		err = errors.New("Error: no password was provided.")
		return
	}

	hpass, err := bEncrypt(password)
	if err != nil {
		return
	}
	u.Password = string(hpass)

	return
}

func bEncrypt(pass string) (string, error) {
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	return string(hashedpass), err
}

//If a login checks out, return a new user object
func (login *Login) Login() (user *User, err error) {
	user, err = FindOneByHandle(login.Handle)
	//If a valid user couldn't be found, the user shall become logged out
	if err != nil {
		user = new(User)
		err = errors.New("The username or password was invalid")
		return
	}

	//If the user account's password validates, update u with the new user object
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.PlaintextPassword))
	if err != nil {
		user = new(User)
		err = errors.New("The username or password was invalid")
		return
	}

	return
}

func FindOneByHandle(handle string) (user *User, err error) {
	//Initialize empty user
	user = new(User)

	//Find one or zero existent users
	FindUserStmt, err := db.Prepare(queries.FindOneByHandle)
	if err != nil {
		return
	}
	defer FindUserStmt.Close()

	//Read any found values into the user object
	err = FindUserStmt.QueryRow(handle).Scan(&user.Id, &user.Handle, &user.Password, &user.Created)
	if err != nil {
		user = new(User)
	}

	//If there are no records, the error will be non-nil
	//You may want to obscure this error message elsewhere, 
	//E.g., to hide whether a login failure was due to a non-existent
	//user or a bad password for a real user.  
	return
}

//If a user with this id exists, return a *User object. Otherwise, nil.
func FindOneUserById(id int64) (user *User, err error) {
	return &User{}, nil

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
