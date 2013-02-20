package main

import (
    "fmt"
    "time"
    "github.com/carbocation/forum.git/forum"
)

func main() {
    makeEntries()
}

func makeEntries() {
    // Make 10 entries based on a skeleton; the Id's will be appropriately distinct.
    N := 10
    skeleton := forum.Entry{ Id: 1, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1}
    var entries = make([]forum.Entry,N)
    for i := 0; i < N; i++ {
        entries[i], entries[i].Id, entries[i].Created = skeleton, int64(i+1), time.Now()
    }

    fmt.Printf("Created new entries with the following data:\n%#v\n",entries)
}