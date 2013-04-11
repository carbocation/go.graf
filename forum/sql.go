package forum

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/carbocation/util.git/datatypes/closuretable"
	_ "github.com/lib/pq"
)

// Retrieves all entries that are descendants of the ancestral entry, including the ancestral entry itself
func RetrieveDescendantEntries(ancestorId string, db *sql.DB) (int64, map[int64]Entry, error) {

	//If the thread ID is not parseable as an integer, stop immediately
	root, err := strconv.ParseInt(ancestorId, 10, 64)
	if err != nil {
		// Default to the root if they gave us a non-integer value
		return 0, map[int64]Entry{}, err
	}

	q := `SELECT e.*
		FROM entry_closures closure
		JOIN entry e ON e.id = closure.descendant
		WHERE closure.ancestor = $1`

	stmt, err := db.Prepare(q)
	defer stmt.Close()
	if err != nil {
		return 0, map[int64]Entry{}, err
	}

	// Query from that prepared statement
	rows, err := stmt.Query(root)
	if err != nil {
		return 0, map[int64]Entry{}, err
	}

	entries := map[int64]Entry{}

	var id, authorid int64
	var title, body string
	var created time.Time

	// Iterate over the rows
	for rows.Next() {
		err = rows.Scan(&id, &title, &body, &created, &authorid)
		if err != nil {
			return 0, map[int64]Entry{}, err
		}

		entries[id] = Entry{Id: id, Title: title, Body: body, Created: created, AuthorId: authorid}
	}

	return root, entries, nil
}

//Returns a closure table built from a given ID
func ClosureTable(id int64) (ct *closuretable.ClosureTable, err error) {
	ct = closuretable.New(id)

	// Pull down the remaining elements in the closure table that are descendants of this node
	q := `select * 
from entry_closures
where descendant in (
select descendant
from entry_closures
where ancestor=$1
)
and ancestor in (
select descendant
from entry_closures
where ancestor=$1
)
and depth = 1`
	stmt, err := Config.DB.Prepare(q)
	if err != nil {
		//fmt.Printf("Statement Preparation Error: %s", err)
		return ct, err
	}

	rows, err := stmt.Query(id)
	if err != nil {
		//fmt.Printf("Query Error: %v", err)
		return ct, err
	}

	//Populate the closuretable
	for rows.Next() {
		var ancestor, descendant int64
		var depth int
		err = rows.Scan(&ancestor, &descendant, &depth)
		if err != nil {
			//fmt.Printf("Rowscan error: %s", err)
			return ct, err
		}

		err = ct.AddChild(closuretable.Child{Parent: ancestor, Child: descendant})

		//err = ct.AddRelationship(closuretable.Relationship{Ancestor: ancestor, Descendant: descendant, Depth: depth})
		if err != nil {
			//fmt.Fprintf(w, "Error: %s", err)
			return ct, err
		}
	}

	return ct, nil
}
