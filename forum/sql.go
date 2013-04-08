package forum

import (
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"time"
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
