/*
This file manages all SQL queries that are made in the forum package.
*/
package forum

var queries = struct {
	DescendantEntries      string
	DescendantClosureTable string
}{
	DescendantEntries: `SELECT e.*
FROM entry_closures closure
JOIN entry e ON e.id = closure.descendant
WHERE closure.ancestor = $1`,
	DescendantClosureTable: `select * 
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
and depth = 1`,
}
