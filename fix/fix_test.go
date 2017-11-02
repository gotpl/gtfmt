package fix

import "testing"

func TestFix(t *testing.T) {
	tpl := `{{index "foobar" "b"}}`
	out, err := Fix(tpl, "index", "strings.Index")
	if err != nil {
		t.Fatal(err)
	}
	expected := `{{strings.Index "foobar" "b"}}`
	if out != expected {
		t.Fatalf("expected:\n%q\n\nbut got:\n%q", expected, out)
	}
}
