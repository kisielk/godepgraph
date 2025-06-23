package main

import (
	"fmt"
	"go/build"
)

// graphvizPrinter implements graphPrinter for the DOT / GraphViz diagramming language.
type graphvizPrinter struct {
	ids map[string]string
}

func newGraphvizPrinter() *graphvizPrinter {
	p := new(graphvizPrinter)
	p.ids = make(map[string]string)
	return p
}

func (p *graphvizPrinter) writeHeader(hLayout bool) {
	fmt.Println("digraph godep {")
	if hLayout {
		fmt.Println(`rankdir="LR"`)
	}
	fmt.Println(`splines=ortho
nodesep=0.4
ranksep=0.8
node [shape="box",style="rounded,filled"]
edge [arrowsize="0.5"]`)
}

func (p *graphvizPrinter) writeNode(pkgName string, attrs *build.Package) {
	id := p.getId(pkgName)

	var color string
	switch {
	case attrs.Goroot:
		color = "palegreen"
	case len(attrs.CgoFiles) > 0:
		color = "darkgoldenrod1"
	case isVendored(attrs.ImportPath):
		color = "palegoldenrod"
	case hasBuildErrors(attrs):
		color = "red"
	default:
		color = "paleturquoise"
	}

	fmt.Printf("%s [label=\"%s\" color=\"%s\" URL=\"%s\" target=\"_blank\"];\n", id, pkgName, color, pkgDocsURL(pkgName))
}

func (p *graphvizPrinter) writeEdge(u string, v string) {
	uId := p.getId(u)
	vId := p.getId(v)
	fmt.Printf("%s -> %s;\n", uId, vId)
}

func (p *graphvizPrinter) getId(pkgName string) string {
	id, ok := p.ids[pkgName]
	if !ok {
		id = "\"" + pkgName + "\""
		p.ids[pkgName] = id
	}
	return id
}

func (p *graphvizPrinter) writeEnd() {
	fmt.Println("}")
}
