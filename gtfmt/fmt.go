package gtfmt

import (
	"fmt"
	"strings"

	"github.com/gotpl/gtfmt/internal/parse"
)

// Formatted reports whether the text in the given template is correctly formatted.
func Formatted(s string) (bool, error) {
	tree, err := parse.ParseNoFuncs("gtfmt", s, "", "")
	if err != nil {
		return false, err
	}
	return s == tree["gtfmt"].Root.String(), nil
}

// Format formats the code inside your template statements without changing any
// other surrounding text.
func Format(s string) (string, error) {
	tree, err := parse.ParseNoFuncs("gtfmt", s, "", "")
	if err != nil {
		return "", err
	}
	return tree["gtfmt"].Root.String(), nil
}

// Fix replaces orig with repl in tpl. tpl must be a valid go template.  Orig
// must be a valid template function name or . path (e.g. .Foo.Bar).  Paths
// *must* start with a ".".
func Fix(tpl string, orig, repl string) (string, error) {
	tree, err := parse.ParseNoFuncs("gtfmt", tpl, "", "")
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
	s.walk(tree["gtfmt"].Root)
	return tree["gtfmt"].Root.String(), nil
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
