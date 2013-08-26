package graf

import (
	"database/sql"
	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var (
	Config         *ConfigFile
	db             *sql.DB                   //db maintains a pool of connections to our database of choice
	store          *sessions.FilesystemStore //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router         *mux.Router               //Dynamic content is managed by graf.Handlers pointed at by the router
	LogWriter      *log.Logger
	ErrorLogWriter *log.Logger
	decoder        *schema.Decoder
)

func Initialize(cfg *ConfigFile, d *sql.DB, s *sessions.FilesystemStore, r *mux.Router, de *schema.Decoder) {
	Config, db, store, router, decoder = cfg, d, s, r, de

	LogWriter = log.New(Config.App.LogAccess, "", 0)
	ErrorLogWriter = log.New(Config.App.LogError, "", 0)
}
