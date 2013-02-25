package forum

import (
	"fmt"
	"github.com/carbocation/util.git/datatypes/binarytree"
	"math/rand"
	"testing"
	"time"
)

func TestClosureConversion(t *testing.T) {
	// Make some sample entries based on a skeleton; the Id's will be appropriately distinct.
	entries := []Entry{
		Entry{Id: 3905, Title: "Hello, world title.", Body: "This is a basic body.", Created: time.Now(), AuthorId: 1},
		Entry{Id: 3906, Title: "Frost Psot", Body: "This is a spam post.", Created: time.Now(), AuthorId: 2},
		Entry{Id: 3907, Title: "Third post", Body: "I want to bury the spam.", Created: time.Now(), AuthorId: 3},
		Entry{Id: 3908, Title: "Les Mis", Body: "It's being shown on the Oscars now.", Created: time.Now(), AuthorId: 3},
		Entry{Id: 3909, Title: "LOOL", Body: "Why are you watching those?", Created: time.Now(), AuthorId: 2},
		Entry{Id: 3910, Title: "Too bad", Body: "I'm here to resurrect the spam.", Created: time.Now(), AuthorId: 2},
	}

	// Create a closure table to represent the relationships among the entries
	// In reality, you'd probably directly import the closure table data into the ClosureTable class
	closuretable := ClosureTable{Relationship{Ancestor: 3905, Descendant: 3905, Depth: 0}}
	closuretable.AddChild(Child{Parent: 3905, Child: 3906})
	closuretable.AddChild(Child{Parent: 3905, Child: 3907})
	closuretable.AddChild(Child{Parent: 3907, Child: 3908})
	closuretable.AddChild(Child{Parent: 3908, Child: 3909})
	closuretable.AddChild(Child{Parent: 3905, Child: 3910})

	//Build a tree out of the entries based on the closure table's instructions.
	tree := walkBody(closuretable.TableToTree(entries))
	expected := "This is a basic body.This is a spam post.I want to bury the spam.It's being shown on the Oscars now.Why are you watching those?I'm here to resurrect the spam."

	if tree != expected {
		t.Errorf("walkBody(tree) yielded %s, expected %s. Have you made a change that caused the iteration order to become indeterminate, e.g., using a map instead of a slice?", tree, expected)
	}
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
