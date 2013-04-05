/*
Derived from zeebo's https://github.com/zeebo/gostbook/blob/master/template.go
*/
package main

import (
	"fmt"
	"html/template"
	"path/filepath"
	"sync"
)

var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.RWMutex

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

func T(name string) *template.Template {
	// First, read from the global map if there is a cached version available:
	cachedMutex.RLock()
	if t, ok := cachedTemplates[name]; ok {
		cachedMutex.RUnlock()
		return t
	}

	// There is no cached version available. Remove the read lock and get a full RW lock, 
	// compile the template, and return it
	cachedMutex.RUnlock()
	cachedMutex.Lock()
	defer cachedMutex.Unlock()

	t := template.New("_base.html").Funcs(funcs)

	t = template.Must(t.ParseFiles(
		"templates/_base.html",
		filepath.Join("templates", name),
	))
	cachedTemplates[name] = t

	return t
}
