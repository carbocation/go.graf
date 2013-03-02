package main

import (
	"fmt"
	"github.com/carbocation/forum.git/forum"
	"github.com/carbocation/util.git/datatypes/binarytree"
	//"math/rand"
	"time"
    "net/http"
)

func main() {
	http.HandleFunc("/hello/", helloHandler)
	http.ListenAndServe("localhost:9999", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	remPartOfURL := r.URL.Path[len("/hello/"):] //get everything after the /hello/ part of the URL
	fmt.Fprintf(w, "Hello %s!", remPartOfURL)

    L := binaryTree()

    fmt.Fprintf(w, "%v", L)
}

func binaryTree() *binarytree.Tree {
	L := binarytree.New(forum.Entry{Id: 1, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
	L.PushLeft(forum.Entry{Id: 2, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
	L.Left().PushLeft(forum.Entry{Id: 3, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})
	L.Left().PushRight(forum.Entry{Id: 4, Title: "Hello, world", Body: "This is a body.", Created: time.Now(), AuthorId: 1})

	channel := binarytree.Walker(L)

	for i := range channel {
		fmt.Printf("%#v\n", i)
	}

    return L
}
