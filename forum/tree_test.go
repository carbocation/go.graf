/*
This is a test of the binary tree / recursive template parsing example given by Rob Pike 
in the go-nuts mailing list. It has been updated to be compatible with go 1.0.3+, including 
1.1beta2. 
*/
package forum

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

// Tree is a binary tree.
type Tree struct {
	Val         int
	Left, Right *Tree
}

// treeTemplate is the definition of a template to walk the tree. We construct it as a template set so it's easy to refer to itself.
// (We could instead drop the {{define "tree"}}...{{end}} and use ExecuteInSet, but sets make all this simpler to set up.)
const treeTemplate = `
        {{define "tree"}}
        [
                {{.Val}}
                {{with .Left}}
                        {{template "tree" .}}
                {{end}}
                {{with .Right}}
                        {{template "tree" .}}
                {{end}}
        ]
        {{end}}
`

func TestTree(t *testing.T) {
	// Build a tree to print.
	var tree = &Tree{
		1,
		&Tree{
			2, &Tree{
				3,
				&Tree{
					4, nil, nil,
				},
				nil,
			},
			&Tree{
				5,
				&Tree{
					6, nil, nil,
				},
				nil,
			},
		},
		&Tree{
			7,
			&Tree{
				8,
				&Tree{
					9, nil, nil,
				},
				nil,
			},
			&Tree{
				10,
				&Tree{
					11, nil, nil,
				},
				nil,
			},
		},
	}
	// Build and parse the set (of one element).
	set := template.Must(template.New("tree").Parse(treeTemplate)) // ALWAYS CHECK ERRORS!!!!!

	var b bytes.Buffer
	// Use set.Execute, starting with the "tree" template.
	// To see the output directly, instead of &b use os.Stdout.
	err := set.Execute(&b, tree) // ALWAYS CHECK ERRORS!!!!!
	if err != nil {
		t.Fatal("exec error:", err)
	}
	// This hoo-hah is to make the comparison easy and clear.
	stripSpace := func(r rune) rune {
		if r == '\t' || r == ' ' || r == '\n' {
			return -1
		}
		return r
	}
	result := strings.Map(stripSpace, b.String())
	const expect = "[1[2[3[4]][5[6]]][7[8[9]][10[11]]]]"
	if result != expect {
		t.Errorf("expected %q got %q", expect, result)
	}
}
