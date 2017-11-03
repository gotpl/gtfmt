# gtfmt

Like go fmt but for go templates.

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