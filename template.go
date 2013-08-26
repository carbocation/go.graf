/*
Derived from zeebo's https://github.com/zeebo/gostbook/blob/master/template.go
*/
package graf

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/carbocation/gotogether"
	"github.com/dustin/go-humanize"
	"github.com/russross/blackfriday"
)

// Note that we can't just preload and cache all of the templates
// because some rely on base templates while others do not
var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.RWMutex

//reverse builds a URL based on route information and paramaters with their arguments
func reverse(name string, things ...interface{}) string {
	//convert the things to strings
	strs := make([]string, len(things))
	for i, th := range things {
		strs[i] = fmt.Sprint(th)
	}
	//grab the route
	u, err := router.GetRoute(name).URL(strs...)
	if err != nil {
		panic(err)
	}
	return u.Path
}

func safeHTML(input string) template.HTML {
	return template.HTML(input)
}

func substring(input string, dropafter int) string {
	if dropafter > len(input) {
		return input
	} else {
		return input[:dropafter]
	}
}

func urlHost(input string) string {
	URL, _ := url.Parse(input)

	return URL.Host
}

func humanizeTime(input time.Time) string {
	return humanize.Time(input)
}

func markDown(input string) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(input)))
}

// From Russ Cox on the go-nuts mailing list
// Modified to treat int and int64 equally, as well as float32 and float64
// https://groups.google.com/d/msg/golang-nuts/OEdSDgEC7js/iyhU9DW_IKcJ
// eq reports whether the first argument is equal to
// any of the remaining arguments.
func eq(args ...interface{}) bool {
	if len(args) == 0 {
		return false
	}
	x := args[0]
	switch x := x.(type) {
	case int:
		for _, y := range args[1:] {
			switch y := y.(type) {
			case int:
				if int64(x) == int64(y) {
					return true
				}
			case int64:
				if int64(x) == int64(y) {
					return true
				}
			}
		}
		return false

	case int64:
		for _, y := range args[1:] {
			switch y := y.(type) {
			case int:
				if int64(x) == int64(y) {
					return true
				}
			case int64:
				if int64(x) == int64(y) {
					return true
				}
			}
		}
		return false

	case float32:
		for _, y := range args[1:] {
			switch y := y.(type) {
			case float32:
				if float64(x) == float64(y) {
					return true
				}
			case float64:
				if float64(x) == float64(y) {
					return true
				}
			}
		}
		return false

	case float64:
		for _, y := range args[1:] {
			switch y := y.(type) {
			case float32:
				if float64(x) == float64(y) {
					return true
				}
			case float64:
				if float64(x) == float64(y) {
					return true
				}
			}
		}
		return false

	case string, byte:
		for _, y := range args[1:] {
			if x == y {
				return true
			}
		}
		return false
	}

	for _, y := range args[1:] {
		if reflect.DeepEqual(x, y) {
			return true
		}
	}
	return false
}

// From gary_b on the go-nuts mailing list
// https://groups.google.com/d/msg/golang-nuts/yGXyPGnHjJQ/ia-zmmOag8IJ
func mapfn(kvs ...interface{}) (map[string]interface{}, error) {
	if len(kvs)%2 != 0 {
		return nil, errors.New("map requires even number of arguments.")
	}
	m := make(map[string]interface{})
	for i := 0; i < len(kvs); i += 2 {
		s, ok := kvs[i].(string)
		if !ok {
			return nil, errors.New("even args to map must be strings.")
		}
		m[s] = kvs[i+1]
	}
	return m, nil
}

var funcs = template.FuncMap{
	"reverse":      reverse,
	"eq":           eq,
	"mapfn":        mapfn,
	"safeHTML":     safeHTML,
	"urlHost":      urlHost,
	"humanizeTime": humanizeTime,
	"substring":    substring,
	"markDown":     markDown,
}

// Parse a template ('name') against _base.html
func T(name string) *template.Template {
	return t("_base.html", name)
}

// Parse a template ('name') against an arbitrary base template.
// Regardless of the base template in use, the 'name' must be unique.
func t(base, name string) *template.Template {
	// First, read from the global cache if available:
	cachedMutex.RLock()
	if t, ok := cachedTemplates[name]; ok {
		defer cachedMutex.RUnlock()
		return t
	}

	// There is no cached version available. Remove the read lock and get a full RW lock,
	// compile the template, and return it
	cachedMutex.RUnlock()
	cachedMutex.Lock()
	defer cachedMutex.Unlock()

	// Create a template with the given basename and custom functions.
	// Panic if there is any error
	n := template.New(base).Funcs(funcs)
	t := template.Must(gotogether.LoadTemplates(n, filepath.Join("templates", base), filepath.Join("templates", name)))

	// Add the newly compiled template to our global cache
	cachedTemplates[name] = t

	return t
}
