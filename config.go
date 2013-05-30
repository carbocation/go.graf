/*
Define exported config types so they can be overwritten.

Note that for things such as the Public config, an alternative
strategy would be to pass the templates an interface{} so you can
pack it with whatever fields you please.
*/
package asksite

import (
	"io"
)

// A config file type is an object that nests various
// public and other config structures
type ConfigFile struct {
	App    *ConfigApp
	DB     *ConfigDB
	Public *ConfigPublic
}

//App-level settings like HTTP ports and secret keys
type ConfigApp struct {
	Identifier  string    //This distinguishes this app for logging and other purposes
	Environment string    //production, dev, ...
	LogAccess   io.Writer //Log every request
	LogError    io.Writer //Errors
	Port        string
	Secret      string
	RootForumID string //The numeric (but stringified) identifier for your root forum
}

//DB connection config
type ConfigDB struct {
	User     string
	Password string
	DBName   string
	Port     string
	PoolSize int //Should be <= max_connections in /etc/postgresql/(version #)/main/postgresql.conf
}

// Public values that can be passed around into e.g., templates
type ConfigPublic struct {
	Site         string //Site name
	Url          string //Full URL, e.g., http://www.google.com
	ContactEmail string //Webmaster email address
	GACode       string //Google Analytics Code
	GAUrl        string //URL of your site, according to Google Analytics
}
