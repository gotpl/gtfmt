# gtfmt [![Build Status](https://travis-ci.org/gotpl/gtfmt.svg?branch=master)](https://travis-ci.org/gotpl/gtfmt)


Like go fmt but for go templates.

Due to several problems found in round-tripping templates through the stdlib
parser and lexer, use of this command should be considered experimental.  See
the bugs section below.

Rewriting with -r is still experimental and should be used with extreme caution,
as it may have unintended consequences.

## Bugs 

Whitespace deletion(i.e. {{- "foo" -}}) gets formatted into actually deleting
the whitespace.  While this technically doesn't change the output of your
template, it is effectively "destructive" to the template, since you can't just
remove a dash if you don't like the effect of the whitespace deletion.

sub-templates (i.e. `{{define "foo"}`}) are currently not supported
(gtfmt will refuse to run on templates that use subtemplates).

Template comments are stripped when formatting.

## Examples
```
$ echo 'Hi!  {{  foo  .Index.Bar  "byte"  }}33' | gtfmt
Hi!  {{foo .Index.Bar "byte"}}33

// replace a function name
$ echo 'Hi!  {{  Foo  .Index.Foo  "Foo"  }}33' | gtfmt -r 'Foo -> Baz'
Hi!  {{Baz .Index.Foo "Foo"}}33

// replace a field's path
$ echo 'Hi!  {{  Foo  .Index.Foo  "Foo"  }}33' | gtfmt -r '.Index.Foo -> .Index.Baz.Foo'
Hi!  {{Baz .Index.Baz.Foo "Foo"}}33
```

## Usage

```
usage: gtfmt [options] [file1] <[file2]...>

Reformats one or more go templates. If not given a filename, will read from stdin.

Options:
  -l    list templates that would be updated (but don't update them)
  -r string
        rewrite rule e.g. '.Foo.Bar -> .Foo.Baz.Bar'


Rewrite rules:
  ** this is still in alpha and subject to change **
  ** use with extreme caution **

  * Replace all or part of a path with another path:
    .Index.Foo -> .Index.Baz.Foo

    The paths *must* start with a ".".  Matching is case sensitive, on full words only.

  * Replace a function with another function:
    foo -> bar

    The lack of a . indicates this is a function replacement.
```
