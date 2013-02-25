package forum

import (
	"errors"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"sort"
)

// A ClosureTable should represent every direct-line relationship, including self-to-self
type ClosureTable []Relationship

// A Relationship is the fundamental unit of the closure table. A relationship is 
// defined between every entry and itselft, its parent, and any of its parent's ancestors.
type Relationship struct {
	Ancestor   int64
	Descendant int64
	Depth      int
}

// A Child is intended to be an ephemeral entity that gets validated and converted to a Relationship
type Child struct {
	Parent int64
	Child  int64
}

type EmptyTableError int

func (e EmptyTableError) Error() error {
	return errors.New("forum: The closure table is empty, so a parent cannot exist, so a child cannot be added.")
}

func ParentDoesNotExistError() error {
	return errors.New("forum: The closure table contains no record of the requested parent, so no child can be created.")
}

func EntityExistsError() error {
	return errors.New("forum: The entity that you are trying to add to the closure table already exists within it. This operation is not permitted.")
}

// AddChild takes a Child, verifies that it is acceptable, verifies that the 
// ClosureTable is suitable to accept a child, and then creates the appropriate 
// Relationships within the ClosureTable to instantiate that child.
func (table *ClosureTable) AddChild(new Child) error {
	if len(*table) < 1 {
		return EmptyTableError.Error(1)
	}

	if table.EntityExists(new.Parent) != true {
		return ParentDoesNotExistError()
	}

	if table.EntityExists(new.Child) {
		return EntityExistsError()
	}

	// It checks out, create all of the consequent ancestral relationships:
	// Self
	*table = append(*table, Relationship{Ancestor: new.Child, Descendant: new.Child, Depth: 0})

	// All derived relationships, including the direct parent<->child relationship
	for _, rel := range table.GetAncestralRelationships(new.Parent) {
		*table = append(*table, Relationship{Ancestor: rel.Ancestor, Descendant: new.Child, Depth: rel.Depth + 1})
	}

	return nil
}

func (table *ClosureTable) GetAncestralRelationships(id int64) []Relationship {
	list := []Relationship{}
	for _, rel := range *table {
		if rel.Descendant == id {
			list = append(list, rel)
		}
	}

	return list
}

// EntityExists asks if an entity of a given id exists in the closure table
// Entities that exist are guaranteed to appear at least once in ancestor and 
// descendant thanks to the self relationship, so the choice of which one to inspect 
// is arbitrary
func (table *ClosureTable) EntityExists(id int64) bool {
	for _, r := range *table {
		if r.Descendant == id {
			return true
		}
	}

	return false
}

// Return the id of the root node of the closure table.
// This method assumes that there can only be one root node.
func (table *ClosureTable) RootNodeId() (int64, error) {
	m := map[int64]int{}
	for _, rel := range *table {
		//In go, it's valid to increment an integer in a map without first zeroing it
		m[rel.Descendant]++
	}

	trip := 0
	var result int64
	for item, count := range m {
		if count == 1 {
			result = item
			trip++
		}

		if trip > 1 {
			return int64(-1), errors.New("More than one potential root node was present in the closure table.")
		}
	}

	if trip < 1 {
		return int64(-1), errors.New("No potential root nodes were present in the closure table.")
	}

	return result, nil
}

// Takes a map of entries whose keys are the same values as the IDs of the closure table entries
// Returns a well-formed *binarytree.Tree with those entries as values.
func (table *ClosureTable) TableToTree(entries []Entry) *binarytree.Tree {
	// Create the tree from the root node:
	forest := map[int64]*binarytree.Tree{}

	// All entries now are trees
	for _, entry := range entries {
		forest[entry.Id] = binarytree.New(entry)
	}

	childparent := table.DepthOneRelationships()

	for _, rel := range childparent {
		// Add the children.
		// If there is already a child, then traverse right until you find nil
		parentTree := forest[rel.Ancestor]
		siblingMode := false

		for {
			if siblingMode {
				if parentTree.Right() == nil {
					// We found an empty slot
					parentTree.SetRight(forest[rel.Descendant])
					forest[rel.Descendant].SetParent(parentTree)
					break
				} else {
					parentTree = parentTree.Right()
				}
			} else {
				if parentTree.Left() == nil {
					// We found an empty slot
					parentTree.SetLeft(forest[rel.Descendant])
					forest[rel.Descendant].SetParent(parentTree)
					break
				} else {
					parentTree = parentTree.Left()
					siblingMode = true
				}
			}
		}
	}

	rootNodeId, err := table.RootNodeId()
	if err != nil {
		return &binarytree.Tree{}
	}

	return forest[rootNodeId]
}

// Returns a map of the ID of each node along with its maximum depth
func (table *ClosureTable) DeepestRelationships() ([]int, map[int][]Relationship) {
	tmp := map[int64]Relationship{}
	out := map[int][]Relationship{}
	discreteDepths := []int{}

	for _, rel := range *table {
		if rel.Depth > tmp[rel.Descendant].Depth {
			tmp[rel.Descendant] = rel
		}
	}

	for _, rel := range tmp {
		out[rel.Depth] = append(out[rel.Depth], rel)
	}

	for depth, _ := range out {
		discreteDepths = append(discreteDepths, depth)
	}

	sort.Ints(discreteDepths)

	return discreteDepths, out
}

// Returns a map of the ID of each node along with its immediate parent
func (table *ClosureTable) DepthOneRelationships() []Relationship {
	out := []Relationship{}

	for _, rel := range *table {
		if rel.Depth == 1 {
			out = append(out, rel)
		}
	}

	return out
}
