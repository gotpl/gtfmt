package gtfmt

import (
	"testing"
)

func TestFixFunc(t *testing.T) {
	tpl := `{{index "index" "d"}}`
	out, err := Fix(tpl, "index", "strings.Index")
	if err != nil {
		t.Fatal(err)
	}
	expected := `{{strings.Index "index" "d"}}`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}
}

func TestFixPath(t *testing.T) {
	tpl := `{{.Foo.Bar ".Foo.Bar"}}`
	out, err := Fix(tpl, ".Foo.Bar", ".Foo.Baz")
	if err != nil {
		t.Fatal(err)
	}
	expected := `{{.Foo.Baz ".Foo.Bar"}}`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}
}

func TestFormat(t *testing.T) {
	tpl := `{{  index   "index"   "d"  }}`
	out, err := Format(tpl)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{{index "index" "d"}}`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}
}

func TestFixFuncWithPath(t *testing.T) {
	tpl := `Hi!  {{  Foo  .Index.Foo  "Foo"  }}33`
	out, err := Fix(tpl, "Foo", "Bar")
	if err != nil {
		t.Fatal(err)
	}
	expected := `Hi!  {{Bar .Index.Foo "Foo"}}33`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}

}
