package asksite

import (
	"os"
)

type environment string

func Environment() *ConfigFile {
	res := &ConfigFile{
		//These are passed to templates
		Public: &ConfigPublic{
			Site:         "localhost:9999",
			Url:          "http://localhost:9999",
			ContactEmail: "james@askbitcoin.com",
			GACode:       "",
			GAUrl:        "",
		},

		DB: &ConfigDB{
			User:     "askbitcoin",
			Password: "xnkxglie",
			DBName:   "projects",
			Port:     "5432",
			PoolSize: 10,
		},

		App: &ConfigApp{
			identifier:  "askbitcoin",
			Environment: "dev",
			LogAccess:   os.Stdout,
			LogError:    os.Stderr,

			//Port that nginx (for reverse proxy) or the browser has to be pointed at
			Port: "9999",

			//64 bit random string generated with `openssl rand -base64 64`
			Secret: `Qkp7F8uW/D8lXdAHKA5dmFGcsvuZkUQKLtpQcM45rUcjYHO05cG+ohr1zf0DwlRughxOEHVhgNOBtZuo5UGbnA==`,
		},
	}

	return res
}
