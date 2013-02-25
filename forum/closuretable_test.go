package forum

import (
    "testing"
    "time"
    "math/rand"
    "fmt"
)

func TestPopulation(t *testing.T) {
    // Number of entries to create for the test
    N := 200

    //Create randomized entries
    entries := makeEntries(N)

    //Create a randomized closure table to structure the entries
    closuretable := buildClosureTable(N)

    //Build a tree out of the entries based on the closure table's instructions.
    tree := closuretable.TableToTree(entries)
    

    fmt.Println(tree)
}

// Helper function that makes the posts.
func makeEntries(N int) map[int64]Entry {
    // Make 10 entries based on a skeleton; the Id's will be appropriately distinct.
    skeleton := Entry{ Id: 1, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1}
    var entries = make(map[int64]Entry,N)
    for i := 0; i < N; i++ {
        skeleton.Id, skeleton.Created, entries[int64(i)] = int64(i), time.Now(), skeleton
    }

    return entries
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