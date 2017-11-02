package fix

import (
	"fmt"

	"github.com/gotpl/gtfix/internal/parse"
)

// Fix replaces orig with repl in tpl.  tpl must be a valid go template.  Orig
// must be a valid template  FieldNode or IdenitifierNode (i.e. a function or a
// .Foo.Bar.Baz).
func Fix(tpl string, orig, repl string) (string, error) {
	tree, err := parse.Parse("", tpl, "", "", nil)
	if err != nil {
		return "", err
	}
	for k, v := range tree {
		fmt.Printf("%q:%#v", k, v.Root.Nodes)
	}
	return tpl, nil
}
