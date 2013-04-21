/*
This file manages all SQL queries that are made in the forum package.
*/
package forum

var queries = struct {
	DescendantEntries         string //Entry itself and all descendents
	DescendantClosureTable    string
	DepthOneDescendantEntries string //Entry itself and all immediate descendents
	DepthOneClosureTable      string
}{
	DescendantEntries: `SELECT e.id, e.title, e.body, e.url, e.created, e.author_id, e.forum
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
	DepthOneDescendantEntries: `SELECT e.id, e.title, e.body, e.url, e.created, e.author_id, e.forum
FROM entry_closures closure
JOIN entry e ON e.id = closure.descendant
WHERE 1=1
AND closure.ancestor = $1
AND (closure.depth=1 OR closure.depth=0)`,
	DepthOneClosureTable: `select * 
from entry_closures
where ancestor=$1
and depth=1`,
}
