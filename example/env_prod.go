/*
Production settings go here.

Obviously, best practices would be to avoid hardcoding database passwords,
etc, into this file, and instead to pull from environment variables (or the like).
*/
package main

import (
	"fmt"
	"os"

	"github.com/carbocation/go.gtfo"
)

type environment string

func Environment() *asksite.ConfigFile {
	var logdir = `/tmp`

	res := &asksite.ConfigFile{
		//These are passed to templates
		Public: &asksite.ConfigPublic{
			Site:         "GTFO: Golang Threaded Forum, Opensource",
			Url:          "http://example.com",
			ContactEmail: "james@example.com",
			GACode:       "UA-00000000-0",
			GAUrl:        "example.com",
		},

		DB: &asksite.ConfigDB{
			User:     "asksite",
			Password: "test",
			DBName:   "projects",
			Port:     "5432",
			PoolSize: 95,
		},

		App: &asksite.ConfigApp{
			LogAccess:   os.Stdout,
			LogError:    os.Stderr,
			Identifier:  "askgolang",
			Environment: "production",

			//Port that nginx (for reverse proxy) or the browser has to be pointed at
			Port: "9996",

			//64 bit random string generated with `openssl rand -base64 64`
			Secret: `/TsvkZlJD/ZLtx+ffq4ldgupCneonDNUmCp8jpXx4ECqRX9LF5JoI9BWH5ysBtjjUcAsLyEwHNZ8X360jBP+tw==`,

			//The ID
			RootForumID: "49",
		},
	}

	al, err := os.OpenFile(fmt.Sprintf(logdir+`/%s_access.log`, res.App.Identifier), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	res.App.LogAccess = al

	el, err := os.OpenFile(fmt.Sprintf(logdir+`/%s_error.log`, res.App.Identifier), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	res.App.LogError = el

	return res
}
