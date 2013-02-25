package forum

import (
    "testing"
    "time"
    "math/rand"
    "fmt"
    "github.com/carbocation/util.git/datatypes/binarytree"
)

func TestPopulation(t *testing.T) {
    // Make some sample entries based on a skeleton; the Id's will be appropriately distinct.
    entries := map[int64]Entry{
        0: Entry{ Id: 0, Title: "Hello, world title.", Body: "This is a basic body.", Created: time.Now(), AuthorId: 1},
        1: Entry{ Id: 1, Title: "Frost Psot", Body: "This is a spam post.", Created: time.Now(), AuthorId: 2},
        2: Entry{ Id: 2, Title: "Third post", Body: "I want to bury the spam.", Created: time.Now(), AuthorId: 3},
        3: Entry{ Id: 3, Title: "Les Mis", Body: "It's being shown on the Oscars now.", Created: time.Now(), AuthorId: 3},
        4: Entry{ Id: 4, Title: "LOOL", Body: "Why are you watching those?", Created: time.Now(), AuthorId: 2},
    }

    // Create a closure table to represent the relationships among the entries
    // In reality, you'd probably directly import the closure table data into the ClosureTable class
    closuretable := ClosureTable{Relationship{Ancestor: 0, Descendant: 0, Depth: 0}}
    closuretable.AddChild(Child{Parent: 0, Child: 1})
    closuretable.AddChild(Child{Parent: 0, Child: 2})
    closuretable.AddChild(Child{Parent: 2, Child: 3})
    closuretable.AddChild(Child{Parent: 3, Child: 4})

    //Build a tree out of the entries based on the closure table's instructions.
    tree := closuretable.TableToTree(entries)
    
    fmt.Println(walkBody(tree))
}

func walkBody(el *binarytree.Tree) string {
    if el == nil {
        return ""
    }

    out := ""
    out += el.Value.(Entry).Body
    out += walkBody(el.Left())
    out += walkBody(el.Right())

    return out
}

func buildClosureTable(N int) ClosureTable {
    // Create the closure table with a single progenitor
    ct := ClosureTable{Relationship{Ancestor: 0, Descendant: 0, Depth: 0}}
    
    for i := 1; i < N; i++ {
        // Create a place for entry #i, making it the child of a random entry j<i
        err := ct.AddChild(Child{Parent: rand.Int63n(int64(i)), Child: int64(i)})
        if err != nil {
            fmt.Println(err)
            break
        }
    }

    return ct
}