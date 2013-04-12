/*
Model keeps track of converters and manipulators of models.

Models and their methods are generally kept in their own 
files (e.g., user.go). This file exists for glue code that 
allows models to be used with libraries such as the Gorilla 
Web Toolkit's 'Schema' library.
*/
package main

import (
	"reflect"

	"github.com/gorilla/schema"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.RegisterConverter(Password(""), convertPassword)
}

func convertPassword(value string) reflect.Value {
	return reflect.ValueOf(Password(value))
}