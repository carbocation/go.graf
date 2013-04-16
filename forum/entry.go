package forum

import (
	"time"

	"github.com/carbocation/util.git/datatypes/closuretable"
)

// Put ModifiedBy, ModifiedAuthor in a separate table. A post can only be 
// created once but modified an infinite number of times.
type Entry struct {
	Id       int64     "The ID of the post"
	Title    string    "Title of the post. Will be empty for entries that are really intended to be comments."
	Body     string    "Contents of the post. Will be empty for entries that are intended to be links."
	Created  time.Time "Time at which the post was created."
	AuthorId int64     "ID of the author of the post"
}

//Note: why not just manage what is real and what is not through methods?
//The 'get' methods would populate the view-type methods automatically
//or leave them blank if irrelevant. The 'set' methods would store the 
//essential / non-derived fields and ignore the others. That way you don't
//have to juggle two view types
type EntryView struct {
	Entry
	Points   int64
	HasVoted bool
}

// Retrieves all entries that are descendants of the ancestral entry, including the ancestral entry itself
func DescendantEntries(root int64) (entries map[int64]Entry, err error) {
	stmt, err := Config.DB.Prepare(queries.DescendantEntries)
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
	stmt, err := Config.DB.Prepare(queries.DescendantClosureTable)
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
