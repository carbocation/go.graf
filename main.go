package main

import (
    "fmt"
    "time"
    "github.com/carbocation/forum.git/forum"
    "math/rand"
    "github.com/carbocation/util.git/datatypes/binarytree"
)

var N int = 4

func main() {
    //makeEntries()
    //closureTable()
    //closureTree()
    //binaryTree()
    fmt.Println("Welcome")
}

func binaryTree() {
    L := binarytree.New(forum.Entry{ Id: 1, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
    L.PushLeft(forum.Entry{ Id: 2, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
    L.Left().PushLeft(forum.Entry{ Id: 3, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
    L.Left().PushRight(forum.Entry{ Id: 4, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
    
    channel := binarytree.Walker(L)

    for i := range channel {
        fmt.Printf("%#v\n", i)
    }
}

func closureTree() {
    table := buildClosureTable(N)
    
    //Now we have a closure table. Feed it to the ClosureTree to build a recursive structure.
    tree := forum.ClosureTree{}
    err := tree.Populate(table)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Printf("%v\n", table)
    fmt.Printf("%v\n", tree)
}

func closureTable() {
    // Create the closure table with a single progenitor
    table := buildClosureTable(N)

    fmt.Printf("%v\n", table)
}

func makeEntries() {
    // Make 10 entries based on a skeleton; the Id's will be appropriately distinct.
    skeleton := forum.Entry{ Id: 1, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1}
    var entries = make([]forum.Entry,N)
    for i := 0; i < N; i++ {
        entries[i], entries[i].Id, entries[i].Created = skeleton, int64(i+1), time.Now()
    }

    fmt.Printf("Created new entries with the following data:\n%#v\n",entries)
}


func buildClosureTable(N int) forum.ClosureTable {
    // Create the closure table with a single progenitor
    ct := forum.ClosureTable{forum.Relationship{Ancestor: 0, Descendant: 0, Depth: 0}}
    
    for i := 1; i < N; i++ {
        // Create a place for entry #i, making it the child of a random entry j<i
        err := ct.AddChild(forum.Child{Parent: rand.Int63n(int64(i)), Child: int64(i)})
        if err != nil {
            fmt.Println(err)
            break
        }
    }

    return ct
}