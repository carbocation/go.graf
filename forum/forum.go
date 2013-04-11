package forum

import (
	"database/sql"
)

type conf struct {
	DB *sql.DB "A live database object"
}

//Create a package-global config object holding needed globals
var Config *conf = &conf{}

//Niladic function to setup the forum
func CreateWith(db *sql.DB) {
	Config.DB = db
}