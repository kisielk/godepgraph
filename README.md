godepgrah
=========

godepgraph is a program for generating a dependency graph of Go packages.

Install
-------

    go get github.com/kisielk/godepgraph
    go install github.com/kisielk/godepgraph

Use
---

For basic usage, just give the package path of interest as the first argument:

    godepgraph github.com/kisielk/godepgraph

By default godepgraph will display packages in the standard library in the graph, though it will not delve in to their dependencies.

If you want to ignore standard library packages entirely, use the -s flag:

    godepgraph -s github.com/kisielk/godepgraph

