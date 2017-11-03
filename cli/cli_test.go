package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseReplace(t *testing.T) {
	stdout := &bytes.Buffer{}
	c, err := Parse(stdout, []string{"-r", "index -> strings.Index"})
	if err != nil {
		t.Fatal(err)
	}
	if s := stdout.String(); s != "" {
		t.Fatalf("expected no stdout, but got %q", s)
	}
	expected := &Command{
		Orig:    "index",
		Replace: "strings.Index",
		Files:   []string{},
	}
	if !reflect.DeepEqual(expected, c) {
		t.Fatalf("Expected:\n%#v\n\ngot:\n%#v", expected, c)
	}
}

func TestReplaceStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := bytes.NewBufferString(`{{index "index" "d"}}`)
	code := ParseAndRun(&stdout, &stderr, stdin, []string{"-r", "index -> strings.Index"})
	if code != 0 {
		t.Errorf("expected code 0 but got %d", code)
	}
	expected := `{{strings.Index "index" "d"}}`

	if s := stdout.String(); s != expected {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, s)
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Expected no stderr but got %q", s)
	}
}

func TestFmtStdin(t *testing.T) {
	var stdout, stderr bytes.Buffer
	stdin := bytes.NewBufferString(`{{  index   "index"   "d"  }}`)
	code := ParseAndRun(&stdout, &stderr, stdin, []string{})
	if code != 0 {
		t.Errorf("expected code 0 but got %d", code)
	}
	expected := `{{index "index" "d"}}`

	if s := stdout.String(); s != expected {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, s)
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Expected no stderr but got %q", s)
	}
}

func TestFmtFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	f1 := filepath.Join(dir, "foo")
	f2 := filepath.Join(dir, "foo2")
	err = ioutil.WriteFile(f1, []byte(`{{  index   "index"   "d"  }}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(f2, []byte(`{{  foo   "bar"   "d"  }}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := ParseAndRun(&stdout, &stderr, nil, []string{f1, f2})
	if code != 0 {
		t.Errorf("expected code 0 but got %d", code)
	}
	expected := []byte(`{{index "index" "d"}}`)
	b, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, b)
	}
	expected = []byte(`{{foo "bar" "d"}}`)
	b, err = ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, b)
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Expected no stderr but got %q", s)
	}
	if s := stdout.String(); s != "" {
		t.Errorf("Expected no stdout but got %q", s)
	}
}

func TestReplaceFiles(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	f1 := filepath.Join(dir, "foo")
	f2 := filepath.Join(dir, "foo2")
	err = ioutil.WriteFile(f1, []byte(`{{  index   "index"   "d"  }}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(f2, []byte(`{{  foo   "bar"   "d"  }}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := ParseAndRun(&stdout, &stderr, nil, []string{"-r", "index -> strings.Index", f1, f2})
	if code != 0 {
		t.Errorf("expected code 0 but got %d", code)
	}
	expected := []byte(`{{strings.Index "index" "d"}}`)
	b, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, b)
	}
	expected = []byte(`{{foo "bar" "d"}}`)
	b, err = ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, expected) {
		t.Errorf("expected:\n%s\nbut got:\n%s", expected, b)
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Expected no stderr but got %q", s)
	}
	if s := stdout.String(); s != "" {
		t.Errorf("Expected no stdout but got %q", s)
	}
}

func TestList(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	f1 := filepath.Join(dir, "foo")
	f2 := filepath.Join(dir, "foo2")
	cont1 := []byte(`{{  index   "index"   "d"  }}`)
	err = ioutil.WriteFile(f1, cont1, 0600)
	if err != nil {
		t.Fatal(err)
	}
	cont2 := []byte(`{{foo "bar" "d"}}`)
	err = ioutil.WriteFile(f2, cont2, 0600)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := ParseAndRun(&stdout, &stderr, nil, []string{"-l", f1, f2})
	if code != 0 {
		t.Errorf("expected code 0 but got %d", code)
	}
	b, err := ioutil.ReadFile(f1)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, cont1) {
		t.Error("contents of unformatted file were changed but should not have been")
	}
	b, err = ioutil.ReadFile(f2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(b, cont2) {
		t.Error("contents of formatted file were changed but should not have been")
	}
	if s := stderr.String(); s != "" {
		t.Errorf("Expected no stderr but got %q", s)
	}
	expected := f1 + "\n"
	if s := stdout.String(); s != expected {
		t.Errorf("Expected only file1 to be listed, but got\n%s", s)
	}
}
