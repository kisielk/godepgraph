# godepgrah

godepgraph is a program for generating a dependency graph of Go packages.

## Install

    go get github.com/kisielk/godepgraph


## Use

For basic usage, just give the package path of interest as the first
argument:

    godepgraph github.com/kisielk/godepgraph

The output is a graph in [Graphviz][graphviz] dot format. If you have the
graphviz tools installed you can render it by piping the output to dot:

    godepgraph github.com/kisielk/godepgraph | dot -Tpng -o godepgraph.png

By default godepgraph will display packages in the standard library in the
graph, though it will not delve in to their dependencies.

## Colors

godepgraph uses a simple color scheme to denote different types of packages:

  * *green*: a package that is part of the Go standard library, installed in `$GOROOT`.
  * *blue*: a regular Go package found in `$GOPATH`.
  * *orange*: a package found in `$GOPATH` that uses cgo by importing the special package "C".

## Ignoring Imports

### The Go Standard Library

If you want to ignore standard library packages entirely, use the -s flag:

    godepgraph -s github.com/kisielk/godepgraph

### By Name

Import paths can be included in a comma-separated list passed to the -i flag:

    godepgraph -i github.com/foo/bar,github.com/baz/blah github.com/something/else

The packages and their imports will be excluded from the graph, unless the imports
are also imported by another package which is not excluded.

### By Prefix

Import paths can also be ignored by prefix. The -p flag takes a comma-separated
list of prefixes:

    godepgraph -p github.com,launchpad.net bitbucket.org/foo/bar



Example
-------
Here's some example output for a component of Gary Burd's [gopkgdoc][gopkgdoc] project:

![Example output](example.png)

[graphviz]: http://graphviz.org
[gopkgdoc]: https://github.com/garyburd/gopkgdoc

