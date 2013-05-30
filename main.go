package asksite

import (
	"database/sql"
	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var (
	Config    *ConfigFile
	db        *sql.DB                                     //db maintains a pool of connections to our database of choice
	store     *sessions.FilesystemStore                   //With an empty first argument, this will put session files in os.TempDir() (/tmp)
	router    *mux.Router               = mux.NewRouter() //Dynamic content is managed by asksite.Handlers pointed at by the router
	LogWriter *log.Logger
)

func Initialize(cfg *ConfigFile, d *sql.DB, s *sessions.FilesystemStore, r *mux.Router) {
	Config, db, store, router = cfg, d, s, r

	LogWriter = log.New(Config.App.LogAccess, "", 0)
}
