/*
Model keeps track of converters and manipulators of models.

Models and their methods are generally kept in their own 
files (e.g., user.go). This file exists for glue code that 
allows models to be used with libraries such as the Gorilla 
Web Toolkit's 'Schema' library.
*/
package asksite

import (
	"reflect"

	"github.com/gorilla/schema"
)

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.RegisterConverter(string(""), convertPassword)
}

func convertPassword(value string) reflect.Value {
	return reflect.ValueOf(string(value))
}