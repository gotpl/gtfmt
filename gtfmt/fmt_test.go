package gtfmt

import (
	"testing"
)

func TestFixFunc(t *testing.T) {
	tpl := `{{index "index" "d"}}`
	out, err := Fix("tpl", tpl, "index", "strings.Index")
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
	out, err := Fix("tpl", tpl, ".Foo.Bar", ".Foo.Baz")
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
	out, err := Format("tpl", tpl)
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
	out, err := Fix("tpl", tpl, "Foo", "Bar")
	if err != nil {
		t.Fatal(err)
	}
	expected := `Hi!  {{Bar .Index.Foo "Foo"}}33`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}

}

func TestNoParseSubTemplate(t *testing.T) {
	tpl := `
{{define "foo" }}
{{ bar 1 }}
{{end}}
{{ template "foo" . }}
`
	_, err := Format("tpl", tpl)
	if err == nil {
		t.Fatal("expected subtemplate to trigger error")
	}
	if err.Error() != "sub templates not currently supported" {
		t.Fatalf("wrong error message from subtemplate: %v", err)
	}
}
