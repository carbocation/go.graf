package forum

import (
	"time"
	"database/sql"

	"github.com/carbocation/go.util/datatypes/closuretable"
)

// Put ModifiedBy, ModifiedAuthor in a separate table. A post can only be 
// created once but modified an infinite number of times.
type Entry struct {
	Id       int64     "The ID of the post"
	Title    string    "Title of the post. Will be empty for entries that are really intended to be comments."
	Body     string    "Contents of the post. Will be empty for entries that are intended to be links."
	Url      string    //Used if the post is just a link
	Created  time.Time "Time at which the post was created."
	AuthorId int64     "ID of the author of the post"
	Forum    bool      `schema:"-"` //Is this Entry actually a forum?
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
	entries, err = getEntries(root, "AllDescendants")
	return
}

// Retrieves entries that are immediate descendants of the ancestral entry, including the ancestral entry itself
func DepthOneDescendantEntries(root int64) (entries map[int64]Entry, err error) {
	entries, err = getEntries(root, "DepthOneDescendants")
	return
}

func getEntries(root int64, flag string) (entries map[int64]Entry, err error) {
	entries = map[int64]Entry{}
	
	var stmt *sql.Stmt

	switch flag {
		case "AllDescendants":
			stmt, err = Config.DB.Prepare(queries.DescendantEntries)
		case "DepthOneDescendants":
			stmt, err = Config.DB.Prepare(queries.DepthOneDescendantEntries)
	}
	if err != nil {
		return
	}
	defer stmt.Close()

	// Query from that prepared statement
	rows, err := stmt.Query(root)
	if err != nil {
		return
	}

	// Iterate over the rows
	for rows.Next() {
		var e Entry
		var body, url sql.NullString
		err = rows.Scan(&e.Id, &e.Title, &body, &url, &e.Created, &e.AuthorId, &e.Forum)
		if err != nil {
			return
		}
		
		//Only the body or the url will be set; they are mutually exclusive
		if body.Valid {
			e.Body = body.String
		}else if url.Valid {
			e.Url = url.String
		}

		entries[e.Id] = e
	}

	return
}

//Returns a closure table of IDs that are descendants of a given ID (or the ID itself)
func ClosureTable(id int64) (ct *closuretable.ClosureTable, err error) {
	ct, err = getClosureTable(id, "AllDescendants")
	return
}

//Returns a closure table keeping only IDs that are direct descendants of a given ID (or the ID itself)
func DepthOneClosureTable(id int64) (ct *closuretable.ClosureTable, err error) {
	ct, err = getClosureTable(id, "DepthOneDescendants")
	return
}

func getClosureTable(id int64, flag string) (ct *closuretable.ClosureTable, err error) {
	ct = closuretable.New(id)
	
	var stmt *sql.Stmt

	// Pull down the remaining elements in the closure table that are descendants of this node
	switch flag {
		case "AllDescendants":
			stmt, err = Config.DB.Prepare(queries.DescendantClosureTable)
		case "DepthOneDescendants":
			stmt, err = Config.DB.Prepare(queries.DepthOneClosureTable)
	}
	if err != nil {
		return
	}
	defer stmt.Close()
	
	rows, err := stmt.Query(id)
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
