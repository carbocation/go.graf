/*
Derived from zeebo's https://github.com/zeebo/gostbook/blob/master/template.go
*/
package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sync"

	"bitbucket.org/carbocation/nrsc"
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

var funcs = template.FuncMap{
	"reverse": reverse,
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
	t := template.Must(nrsc.LoadTemplates(n, filepath.Join("templates", base), filepath.Join("templates", name)))

	// Add the newly compiled template to our global cache
	cachedTemplates[name] = t

	return t
}
