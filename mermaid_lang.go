package main

import (
	"fmt"
	"go/build"
)

// mermaidPrinter implements graphPrinter for the Mermaid diagramming language.
type mermaidPrinter struct {
	ids    map[string]string
	nextId int
}

func newMermaidPrinter() *mermaidPrinter {
	p := new(mermaidPrinter)
	p.ids = make(map[string]string)
	return p
}

func (p *mermaidPrinter) writeHeader(hLayout bool) {
	if hLayout {
		fmt.Println("flowchart LR")
	} else {
		fmt.Println("flowchart TD")
	}

	fmt.Println()
	fmt.Println("classDef goroot fill:#1D4,color:white")
	fmt.Println("classDef cgofiles fill:#D52,color:white")
	fmt.Println("classDef vendored fill:#D90,color:white")
	fmt.Println("classDef buildErrs fill:#C10,color:white")
}

func (p *mermaidPrinter) writeNode(pkgName string, attrs *build.Package) {
	id := p.getId(pkgName)

	var classname string
	switch {
	case attrs.Goroot:
		classname = "goroot"
	case len(attrs.CgoFiles) > 0:
		classname = "cgofiles"
	case isVendored(attrs.ImportPath):
		classname = "vendored"
	case hasBuildErrors(attrs):
		classname = "buildErrs"
	}

	fmt.Println()
	fmt.Printf("%s[%s]\n", id, pkgName)
	fmt.Printf("click %s href %q\n", id, pkgDocsURL(pkgName))

	if classname != "" {
		fmt.Printf("class %s %s\n", id, classname)
	}
}

func (p *mermaidPrinter) writeEdge(u string, v string) {
	uId := p.getId(u)
	vId := p.getId(v)
	fmt.Printf("%s --> %s\n", uId, vId)
}

func (p *mermaidPrinter) getId(pkgName string) string {
	id, ok := p.ids[pkgName]
	if !ok {
		p.nextId = p.nextId + 1
		id = fmt.Sprintf("%d", p.nextId)
		p.ids[pkgName] = id
	}
	return id
}

func (p *mermaidPrinter) writeEnd() {}
