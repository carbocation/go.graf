package forum

import (
    //"errors"
    "fmt"
)

type ClosureTree struct {
    //Is-a:
    Entry "A closure tree element is also an entry element"
    
    //Has-a:
    Value int64
    Children []ClosureTree
}

func (tree *ClosureTree) Populate(table ClosureTable) error {
    // Find the root node(s), as they are the ones without ancestors.
    // We require that there only be one root node.
    rootId, err := table.RootNodeId()
    if err != nil {
        return err
    }

    tree.Value = rootId
    tree.Children = []ClosureTree{}
    //tree.Children = tree.buildTree(rootId, ct)

    fmt.Println("Tree consists of",tree)

    return nil
}

func (tree *ClosureTree) buildTree(currentNode int64, table ClosureTable) []ClosureTree {
    //return []ClosureTree{}
    //return tree.buildTree(
    return []ClosureTree{}
}


/*
type ClosureTree struct {
    Entries map[int64]Entry
    Children []ClosureTree
}
*/