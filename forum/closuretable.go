package forum

import (
    "errors"
)

// A ClosureTable should represent every direct-line relationship, including self-to-self
type ClosureTable []Relationship

// A Relationship is the fundamental unit of the closure table. A relationship is 
// defined between every entry and itselft, its parent, and any of its parent's ancestors.
type Relationship struct {
    Ancestor int64
    Descendant int64
    Depth int
}

// A Child is intended to be an ephemeral entity that gets validated and converted to a Relationship
type Child struct {
    Parent int64
    Child int64
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
func (ct *ClosureTable) AddChild(new Child) error {
    if len(*ct) < 1 {
        return EmptyTableError.Error(1)
    }

    if ct.EntityExists(new.Parent) != true {
        return ParentDoesNotExistError()
    }

    if ct.EntityExists(new.Child) {
        return EntityExistsError()
    }
    
    // It checks out, create the direct relationship, and all derived relationships:
    // Self
    *ct = append(*ct, Relationship{Ancestor: new.Child, Descendant: new.Child, Depth: 0})
    // Direct relationship
    //*ct = append(*ct, Relationship{Ancestor: new.Parent, Descendant: new.Child, Depth: 1})
    for _, rel := range ct.GetAncestralRelationships(new.Parent) {
        *ct = append(*ct, Relationship{Ancestor: rel.Ancestor, Descendant: new.Child, Depth: rel.Depth+1})
    }

    return nil
}

func (ct *ClosureTable) GetAncestralRelationships(id int64) []Relationship {
    list := []Relationship{}
    for _, rel := range *ct {
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
func (ct *ClosureTable) EntityExists(id int64) bool {
    for _, r := range *ct {
        if r.Descendant == id {
            return true
        }
    }

    return false
}
