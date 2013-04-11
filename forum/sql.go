package forum

import (
	"time"

	"github.com/carbocation/util.git/datatypes/closuretable"
)

// Retrieves all entries that are descendants of the ancestral entry, including the ancestral entry itself
func DescendantEntries(root int64) (entries map[int64]Entry, err error) {
	q := `SELECT e.*
		FROM entry_closures closure
		JOIN entry e ON e.id = closure.descendant
		WHERE closure.ancestor = $1`

	stmt, err := Config.DB.Prepare(q)
	defer stmt.Close()
	if err != nil {
		return
	}

	// Query from that prepared statement
	rows, err := stmt.Query(root)
	if err != nil {
		return
	}

	entries = map[int64]Entry{}

	var id, authorid int64
	var title, body string
	var created time.Time

	// Iterate over the rows
	for rows.Next() {
		err = rows.Scan(&id, &title, &body, &created, &authorid)
		if err != nil {
			return
		}

		entries[id] = Entry{Id: id, Title: title, Body: body, Created: created, AuthorId: authorid}
	}

	return
}

//Returns a closure table of IDs that are descendants of a given ID
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
		return
	}

	rows, err := stmt.Query(id)
	defer stmt.Close()
	if err != nil {
		//fmt.Printf("Query Error: %v", err)
		return
	}

	//Populate the closuretable
	for rows.Next() {
		var ancestor, descendant int64
		var depth int
		err = rows.Scan(&ancestor, &descendant, &depth)
		if err != nil {
			//fmt.Printf("Rowscan error: %s", err)
			return
		}

		err = ct.AddChild(closuretable.Child{Parent: ancestor, Child: descendant})

		//err = ct.AddRelationship(closuretable.Relationship{Ancestor: ancestor, Descendant: descendant, Depth: depth})
		if err != nil {
			//fmt.Fprintf(w, "Error: %s", err)
			return
		}
	}

	return
}
