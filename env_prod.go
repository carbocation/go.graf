package main

import (
	"fmt"
	"os"
)

type environment string

func Environment() *ConfigFile {
	res := &ConfigFile{
		//These are passed to templates
		Public: &ConfigPublic{
			Site:         "Ask Bitcoin",
			Url:          "http://askbitcoin.com",
			ContactEmail: "james@askbitcoin.com",
			GACode:       "UA-36655899-3",
			GAUrl:        "askbitcoin.com",
		},

		DB: &ConfigDB{
			User:     "askbitcoin",
			Password: "xnkxglie",
			DBName:   "projects",
			Port:     "5432",
			PoolSize: 95,
		},

		App: &ConfigApp{
			LogAccess:   os.Stdout,
			LogError:    os.Stderr,
			identifier:  "askbitcoin",
			Environment: "production",
			//Port that nginx (for reverse proxy) or the browser has to be pointed at
			Port: "9999",
			//64 bit random string generated with `openssl rand -base64 64`
			Secret: `75Oop7MSN88WstKJSTyu9ALiO0Nbeckv/4/eDLDJcpXn0Ny1H9PdpzXDqApie77tZ04GFsdHehmzcMkAqh16Dg==`,
		},
	}

	al, err := os.OpenFile(fmt.Sprintf(`/data/bin/%s_access.log`, res.App.Identifier()), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	res.App.LogAccess = al

	el, err := os.OpenFile(fmt.Sprintf(`/data/bin/%s_error.log`, res.App.Identifier()), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		panic(err)
	}
	res.App.LogError = el

	return res
}
