package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
)

var (
	pkgs   map[string]*build.Package
	ids    map[string]int
	nextId int

	ignored = map[string]bool{
		"C": true,
	}
)

func init() {
	pkgs = make(map[string]*build.Package)
	ids = make(map[string]int)
}

func processPackage(root string, pkgName string) error {
	if ignored[pkgName] {
		return nil
	}

	pkg, err := build.Import(pkgName, root, 0)
	if err != nil {
		return fmt.Errorf("failed to import %s: %s", pkgName, err)
	}

	pkgs[pkg.ImportPath] = pkg

	// Don't worry about dependencies for stdlib packages
	if pkg.Goroot {
		return nil
	}

	for _, imp := range pkg.Imports {
		if _, ok := pkgs[imp]; !ok {
			if err := processPackage(root, imp); err != nil {
				return err
			}
		}
	}
	return nil
}

func getId(name string) int {
	id, ok := ids[name]
	if !ok {
		id = nextId
		nextId++
		ids[name] = id
	}
	return id
}

func main() {
	ignoreStdlib := flag.Bool("s", false, "ignore packages in the go standard library")
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need one package name to process")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}
	if err := processPackage(cwd, args[0]); err != nil {
		log.Fatal(err)
	}

	fmt.Println("digraph godep {")
	for pkgName, pkg := range pkgs {
		pkgId := getId(pkgName)

		var color string
		if pkg.Goroot {
			if !*ignoreStdlib {
				color = "palegreen"
			} else {
				continue
			}
		} else {
			color = "paleturquoise"
		}

		fmt.Printf("%d [label=\"%s\" style=\"filled\" color=\"%s\"];\n", pkgId, pkgName, color)

		// Don't render imports from packages in Goroot
		if pkg.Goroot {
			continue
		}

		for _, imp := range pkg.Imports {
			if ignored[imp] {
				continue
			}

			impPkg := pkgs[imp]
			if impPkg.Goroot && *ignoreStdlib {
				continue
			}
			impId := getId(imp)
			fmt.Printf("%d -> %d;\n", pkgId, impId)
		}
	}
	fmt.Println("}")
}
