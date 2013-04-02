/*
Derived from zeebo's https://github.com/zeebo/gostbook/blob/master/template.go
*/
package main

import (
	"html/template"
	"path/filepath"
	"sync"
)

var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.RWMutex

var funcs = template.FuncMap{
	"reverse": reverse,
}

func T(name string) *template.Template {
	cachedMutex.RLock()
	defer cachedMutex.RUnlock()

	if t, ok := cachedTemplates[name]; ok {
		return t
	}

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
