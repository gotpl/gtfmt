package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gotpl/gtfmt/gtfmt"
)

// Main is the main entrypoint for the gtfix binary.
func Main() int {
	return ParseAndRun(os.Stdout, os.Stderr, os.Stdin, os.Args[1:])
}

// ParseAndRun parses the command line, and then runs gtfix.
func ParseAndRun(stdout, stderr io.Writer, stdin io.Reader, args []string) int {
	log := log.New(stderr, "", 0)
	c, err := Parse(stdout, args)
	if err == flag.ErrHelp {
		return 2
	}
	if err != nil {
		log.Println("ERROR: ", err)
		return 1
	}
	c.Stderr = stderr
	c.Stdin = stdin
	c.Stdout = stdout
	if err := c.Run(); err != nil {
		log.Println("ERROR: ", err)
		return 1
	}
	return 0
}

// Parse parses the given args.
func Parse(stdout io.Writer, args []string) (*Command, error) {
	fs := flag.FlagSet{}
	fs.SetOutput(stdout)
	c := &Command{}
	var replace string
	fs.StringVar(&replace, "r", "", "rewrite rule e.g. '.Foo.Bar -> .Foo.Baz.Bar'")
	fs.BoolVar(&c.List, "l", false, "list templates that would be updated (but don't update them)")
	fs.Usage = func() {
		fmt.Fprintln(stdout, `usage: gtfmt [options] [file1] <[file2]...>

Reformats one or more go templates. If not given a filename, will read from stdin.

Options:`)
		fs.PrintDefaults()
		fmt.Fprintln(stdout, `

Rewrite rules:
  ** this is still in alpha and subject to change **
  ** use with extreme caution **

  * Replace all or part of a path with another path:
    .Index.Foo -> .Index.Baz.Foo

    The paths *must* start with a ".".  Matching is case sensitive, on full words only.

  * Replace a function with another function:
    foo -> bar

    The lack of a . indicates this is a function replacement.
`)
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if replace != "" {
		vals := strings.Split(replace, " -> ")
		if len(vals) != 2 || vals[0] == "" {
			return nil, errors.New("rewrite rule must be in the format 'foo -> bar'")
		}
		c.Orig = vals[0]
		c.Replace = vals[1]
	}
	c.Files = fs.Args()
	return c, nil
}

// Command is a Command to run.
type Command struct {
	Orig    string
	Replace string
	List    bool // if true, only list what files need formatting
	Files   []string
	Stdout  io.Writer
	Stdin   io.Reader
	Stderr  io.Writer
}

// Run runs the command
func (c *Command) Run() error {
	if c.Orig == "" {
		return c.format()
	}
	return c.replace()
}

func (c *Command) format() error {
	if len(c.Files) == 0 {
		return c.fmtStdin()
	}
	for _, fn := range c.Files {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		orig := string(b)
		if c.List {
			ok, err := gtfmt.Formatted(fn, orig)
			if err != nil {
				return err
			}
			if !ok {
				io.WriteString(c.Stdout, fn+"\n")
			}
			continue
		}
		s, err := gtfmt.Format(fn, orig)
		if err != nil {
			return err
		}
		if s != orig {
			info, err := os.Stat(fn)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(fn, []byte(s), info.Mode())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Command) fmtStdin() error {
	b, err := ioutil.ReadAll(c.Stdin)
	if err != nil {
		return err
	}
	orig := string(b)
	if c.List {
		ok, err := gtfmt.Formatted("stdin", orig)
		if err != nil {
			return err
		}
		if ok {
			io.WriteString(c.Stdout, "formatted\n")
		} else {
			io.WriteString(c.Stdout, "unformatted\n")
		}
		return nil
	}
	s, err := gtfmt.Format("stdin", orig)
	if err != nil {
		return err
	}
	_, err = io.WriteString(c.Stdout, s)
	return err
}

func (c *Command) replace() error {
	if len(c.Files) == 0 {
		return c.replaceStdin()
	}
	for _, fn := range c.Files {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		tpl := string(b)
		s, err := gtfmt.Fix(fn, tpl, c.Orig, c.Replace)
		if err != nil {
			return err
		}
		if c.List {
			if s != tpl {
				io.WriteString(c.Stdout, fn+"\n")
			}
			continue
		}
		if s != tpl {
			info, err := os.Stat(fn)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(fn, []byte(s), info.Mode())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Command) replaceStdin() error {
	b, err := ioutil.ReadAll(c.Stdin)
	if err != nil {
		return err
	}
	tpl := string(b)
	s, err := gtfmt.Fix("stdin", tpl, c.Orig, c.Replace)
	if err != nil {
		return err
	}
	if c.List {
		if s == tpl {
			io.WriteString(c.Stdout, "unchanged\n")
		} else {
			io.WriteString(c.Stdout, "changed\n")
		}
		return nil
	}
	_, err = io.WriteString(c.Stdout, s)
	return err
}
