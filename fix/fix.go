package fix

import (
	"fmt"
	"strings"

	"github.com/gotpl/gtfix/internal/parse"
)

// Fix replaces orig with repl in tpl.  tpl must be a valid go template.  Orig
// must be a valid template function name or . path (e.g. .Foo.Bar).  Paths *must* start with a ".".
func Fix(tpl string, orig, repl string) (string, error) {
	tree, err := parse.ParseNoFuncs("", tpl, "", "")
	if err != nil {
		return "", err
	}
	s := &state{}
	if strings.HasPrefix(orig, ".") {
		// append a dot at the end to ensure we get full word matching
		if strings.HasSuffix(orig, ".") {
			s.path = orig
		} else {
			s.path = orig + "."
		}
		if strings.HasSuffix(repl, ".") {
			s.repl = repl
		} else {
			s.repl = repl + "."
		}
	} else {
		s.fn = orig
		s.repl = repl
	}
	for k, v := range tree {
		if k != "" {
			fmt.Printf("parsing tree %q\n", k)
		}
		s.walk(v.Root)
	}
	vals := make([]string, 0, len(tree))
	for _, v := range tree {
		vals = append(vals, v.Root.String())
	}
	return strings.Join(vals, ""), nil
}

type state struct {
	fn   string
	path string
	repl string
}

// walk steps through the major pieces of the template structure.
func (s *state) walk(node parse.Node) {
	if node == nil {
		return
	}
	switch node := node.(type) {
	case *parse.ActionNode:
		for _, n := range node.Pipe.Cmds {
			s.walk(n)
		}
	case *parse.IfNode:
		s.walkBranch(node.BranchNode)
	case *parse.RangeNode:
		s.walk(node.Pipe)
		s.walkBranch(node.BranchNode)
	case *parse.WithNode:
		s.walk(node.List)
		s.walk(node.Pipe)
		s.walkBranch(node.BranchNode)
	case *parse.ListNode:
		for _, n := range node.Nodes {
			s.walk(n)
		}
	case *parse.TemplateNode:
		s.walk(node.Pipe)
	case *parse.StringNode, *parse.TextNode:
		// nothing to do
	case *parse.IdentifierNode:
		if node.Ident == s.fn {
			node.Ident = s.repl
		}
	case *parse.CommandNode:
		for _, n := range node.Args {
			s.walk(n)
		}
	case *parse.FieldNode:
		// append a dot at the end to ensure we get full word matching
		ident := "." + strings.Join(node.Ident, ".") + "."
		if !strings.Contains(ident, s.path) {
			return
		}
		val := strings.Trim(strings.Replace(ident, s.path, s.repl, 1), ".")
		node.Ident = strings.Split(val, ".")
	default:
		panic(fmt.Sprintf("unknown node: %T", node))
	}
}

func (s *state) walkBranch(node parse.BranchNode) {
	s.walk(node.Pipe)
	s.walk(node.List)
	s.walk(node.ElseList)
}
