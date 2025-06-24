# godepgraph

godepgraph is a program for generating a dependency graph of Go packages.

[![Go](https://github.com/kisielk/godepgraph/actions/workflows/go.yml/badge.svg)](https://github.com/kisielk/godepgraph/actions/workflows/go.yml)

## Install

    go install github.com/kisielk/godepgraph@latest

## Use

For basic usage, just give the package path of interest as the first
argument:

    godepgraph github.com/kisielk/godepgraph
    
If you intend to graph a go mod project, your package should be passed as a relative path:

    godepgraph ./pkg/api

The default output is a graph in [Graphviz][graphviz] dot format. If you have the
graphviz tools installed you can render it by piping the output to dot:

    godepgraph github.com/kisielk/godepgraph | dot -Tpng -o godepgraph.png

You can also generate [Mermaid](https://mermaid.js.org/) graphs:

    godepgraph -format mermaid github.com/kisielk/godepgraph > graph.mmd

By default godepgraph will display packages in the standard library in the
graph, though it will not delve in to their dependencies.

## Colors

godepgraph uses a simple color scheme to denote different types of packages:

  * *green*: a package that is part of the Go standard library, installed in `$GOROOT`.
  * *blue*: a regular Go package found in `$GOPATH`.
  * *yellow*: a vendored Go package found in `$GOPATH`.
  * *orange*: a package found in `$GOPATH` that uses cgo by importing the special package "C".

## Ignoring Imports

### The Go Standard Library

If you want to ignore standard library packages entirely, use the -s flag:

    godepgraph -s github.com/kisielk/godepgraph

### Vendored Libraries

If you want to ignore vendored packages entirely, use the -novendor flag:

    godepgraph -novendor github.com/something/else

### By Name

Import paths can be ignored in a comma-separated list passed to the -i flag:

    godepgraph -i github.com/foo/bar,github.com/baz/blah github.com/something/else

The packages and their imports will be excluded from the graph, unless the imports
are also imported by another package which is not excluded.

### By Prefix

Import paths can also be ignored by prefix. The -p flag takes a comma-separated
list of prefixes:

    godepgraph -p github.com,launchpad.net bitbucket.org/foo/bar

## Example

Here's some example output for godepgraph itself:

Using `godepgraph -format mermaid github.com/kisielk/godepgraph`:

```mermaid
flowchart TD

classDef goroot fill:#1D4,color:white
classDef cgofiles fill:#D52,color:white
classDef vendored fill:#D90,color:white
classDef buildErrs fill:#C10,color:white

1[flag]
click 1 href "https://godoc.org/flag"
class 1 goroot

2[fmt]
click 2 href "https://godoc.org/fmt"
class 2 goroot

3[github.com/kisielk/godepgraph]
click 3 href "https://godoc.org/github.com/kisielk/godepgraph"
3 --> 1
3 --> 2
3 --> 4
3 --> 5
3 --> 6
3 --> 7
3 --> 8

4[go/build]
click 4 href "https://godoc.org/go/build"
class 4 goroot

5[log]
click 5 href "https://godoc.org/log"
class 5 goroot

6[os]
click 6 href "https://godoc.org/os"
class 6 goroot

7[sort]
click 7 href "https://godoc.org/sort"
class 7 goroot

8[strings]
click 8 href "https://godoc.org/strings"
class 8 goroot
```

Using `godepgraph -format graphviz github.com/kisielk/godepgraph`:

![Example output](example.png)

[graphviz]: http://graphviz.org

Here's the output of `godepgraph -format mermaid -nostdlib github.com/kisielk/errcheck`:

```mermaid
flowchart TD

classDef goroot fill:#1D4,color:white
classDef cgofiles fill:#D52,color:white
classDef vendored fill:#D90,color:white
classDef buildErrs fill:#C10,color:white

1[github.com/kisielk/errcheck]
click 1 href "https://godoc.org/github.com/kisielk/errcheck"
1 --> 2
1 --> 3

2[github.com/kisielk/errcheck/errcheck]
click 2 href "https://godoc.org/github.com/kisielk/errcheck/errcheck"
2 --> 4
2 --> 3

5[golang.org/x/mod/semver]
click 5 href "https://godoc.org/golang.org/x/mod/semver"

6[golang.org/x/sync/errgroup]
click 6 href "https://godoc.org/golang.org/x/sync/errgroup"

4[golang.org/x/tools/go/analysis]
click 4 href "https://godoc.org/golang.org/x/tools/go/analysis"

7[golang.org/x/tools/go/gcexportdata]
click 7 href "https://godoc.org/golang.org/x/tools/go/gcexportdata"
7 --> 8

3[golang.org/x/tools/go/packages]
click 3 href "https://godoc.org/golang.org/x/tools/go/packages"
3 --> 6
3 --> 7
3 --> 9
3 --> 10
3 --> 11

12[golang.org/x/tools/go/types/objectpath]
click 12 href "https://godoc.org/golang.org/x/tools/go/types/objectpath"
12 --> 13
12 --> 11

14[golang.org/x/tools/go/types/typeutil]
click 14 href "https://godoc.org/golang.org/x/tools/go/types/typeutil"
14 --> 15

13[golang.org/x/tools/internal/aliases]
click 13 href "https://godoc.org/golang.org/x/tools/internal/aliases"

16[golang.org/x/tools/internal/event]
click 16 href "https://godoc.org/golang.org/x/tools/internal/event"
16 --> 17
16 --> 18
16 --> 19

17[golang.org/x/tools/internal/event/core]
click 17 href "https://godoc.org/golang.org/x/tools/internal/event/core"
17 --> 18
17 --> 19

18[golang.org/x/tools/internal/event/keys]
click 18 href "https://godoc.org/golang.org/x/tools/internal/event/keys"
18 --> 19

19[golang.org/x/tools/internal/event/label]
click 19 href "https://godoc.org/golang.org/x/tools/internal/event/label"

8[golang.org/x/tools/internal/gcimporter]
click 8 href "https://godoc.org/golang.org/x/tools/internal/gcimporter"
8 --> 12
8 --> 13
8 --> 20
8 --> 11

9[golang.org/x/tools/internal/gocommand]
click 9 href "https://godoc.org/golang.org/x/tools/internal/gocommand"
9 --> 5
9 --> 16
9 --> 18
9 --> 19

10[golang.org/x/tools/internal/packagesinternal]
click 10 href "https://godoc.org/golang.org/x/tools/internal/packagesinternal"

20[golang.org/x/tools/internal/pkgbits]
click 20 href "https://godoc.org/golang.org/x/tools/internal/pkgbits"

21[golang.org/x/tools/internal/stdlib]
click 21 href "https://godoc.org/golang.org/x/tools/internal/stdlib"

15[golang.org/x/tools/internal/typeparams]
click 15 href "https://godoc.org/golang.org/x/tools/internal/typeparams"
15 --> 13

11[golang.org/x/tools/internal/typesinternal]
click 11 href "https://godoc.org/golang.org/x/tools/internal/typesinternal"
11 --> 14
11 --> 13
11 --> 21
11 --> 22

22[golang.org/x/tools/internal/versions]
click 22 href "https://godoc.org/golang.org/x/tools/internal/versions"
```
