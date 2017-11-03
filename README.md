# gtfmt [![Build Status](https://travis-ci.org/gotpl/gtfmt.svg?branch=master)](https://travis-ci.org/gotpl/gtfmt)

Like go fmt but for go templates.

Note that "reformatting" only changes the code inside template actions, it will
never change the text outside the template actions.

Reformatting should be considered safe and will not change your template's
output at all.

Rewriting with -r is still experimental and should be used with extreme caution,
as it may have unintended consequences.

## Examples

$ echo 'Hi!  {{  foo  .Index.Bar  "byte"  }}33' | gtfmt
Hi!  {{foo .Index.Bar "byte"}}33

// replace a function name
$ echo 'Hi!  {{  Foo  .Index.Foo  "Foo"  }}33' | gtfmt -r 'Foo -> Baz'
Hi!  {{Baz .Index.Foo "Foo"}}33

// replace a field's path
$ echo 'Hi!  {{  Foo  .Index.Foo  "Foo"  }}33' | gtfmt -r '.Index.Foo -> .Index.Baz.Foo'
Hi!  {{Baz .Index.Baz.Foo "Foo"}}33


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
