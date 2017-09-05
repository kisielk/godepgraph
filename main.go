package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	pkgs   map[string]*build.Package
	ids    map[string]int
	nextId int

	ignored         map[string]bool
	ignoredPrefixes []string
	onlyPrefixes    []string

	ignoreStdlib   = new(bool)
	delveGoroot    = new(bool)
	ignorePrefixes = new(string)
	ignorePackages = new(string)
	onlyPrefix     = new(string)
	tagList        = new(string)
	horizontal     = new(bool)
	includeTests   = new(bool)
	maxLevel       = new(int)

	buildTags    []string
	buildContext           = build.Default
	output       io.Writer = os.Stdout
)

func init() {
	flag.BoolVar(ignoreStdlib, "s", false, "ignore packages in the Go standard library")
	flag.BoolVar(delveGoroot, "d", false, "show dependencies of packages in the Go standard library")
	flag.StringVar(ignorePrefixes, "p", "", "a comma-separated list of prefixes to ignore")
	flag.StringVar(ignorePackages, "i", "", "a comma-separated list of packages to ignore")
	flag.StringVar(onlyPrefix, "o", "", "a comma-separated list of prefixes to include")
	flag.StringVar(tagList, "tags", "", "a comma-separated list of build tags to consider satisified during the build")
	flag.BoolVar(horizontal, "horizontal", false, "lay out the dependency graph horizontally instead of vertically")
	flag.BoolVar(includeTests, "t", false, "include test packages")
	flag.IntVar(maxLevel, "l", 256, "max level of go dependency graph")
}

func main() {
	nextId = 0
	ignored = map[string]bool{
		"C": true,
	}
	pkgs = make(map[string]*build.Package)
	ids = make(map[string]int)
	*ignoreStdlib = false
	*delveGoroot = false
	*ignorePrefixes = ""
	*ignorePackages = ""
	*onlyPrefix = ""
	*tagList = ""
	*horizontal = false
	*includeTests = false
	*maxLevel = 256
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		log.Fatal("need one package name to process")
	}

	if *ignorePrefixes != "" {
		ignoredPrefixes = strings.Split(*ignorePrefixes, ",")
	}
	if *onlyPrefix != "" {
		onlyPrefixes = strings.Split(*onlyPrefix, ",")
	}
	if *ignorePackages != "" {
		for _, p := range strings.Split(*ignorePackages, ",") {
			ignored[p] = true
		}
	}
	if *tagList != "" {
		buildTags = strings.Split(*tagList, ",")
	}
	buildContext.BuildTags = buildTags

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}
	if err := processPackage(cwd, args[0], 0); err != nil {
		log.Fatal(err)
	}

	println("digraph godep {")
	if *horizontal {
		println(`rankdir="LR"`)
	}

	// sort packages
	pkgKeys := []string{}
	for k := range pkgs {
		pkgKeys = append(pkgKeys, k)
	}
	sort.Strings(pkgKeys)

	for _, pkgName := range pkgKeys {
		pkg := pkgs[pkgName]
		pkgId := getId(pkgName)
		pkgUrl := ""
		if isIgnored(pkg) {
			continue
		}

		var color string
		if pkg.Goroot {
			color = "palegreen"
		} else if len(pkg.CgoFiles) > 0 {
			color = "darkgoldenrod1"
		} else {
			pkgUrl = getUrl(pkgName)
			color = "paleturquoise"
		}
		if pkgUrl == "" {
			printf("_%d [label=\"%s\" style=\"filled\" color=\"%s\"];\n", pkgId, pkgName, color)
		} else {
			printf("_%d [label=\"%s\" style=\"filled\" color=\"%s\" URL=\"%s\" target=\"_top\"];\n", pkgId, pkgName, color, pkgUrl)
		}
		// Don't render imports from packages in Goroot
		if pkg.Goroot && !*delveGoroot {
			continue
		}

		for _, imp := range getImports(pkg) {
			impPkg := pkgs[imp]
			if impPkg == nil || isIgnored(impPkg) {
				continue
			}

			impId := getId(imp)
			printf("_%d -> _%d;\n", pkgId, impId)
		}
	}
	println("}")
}

func processPackage(root string, pkgName string, level int) error {
	if level++; level > *maxLevel {
		return nil
	}
	if ignored[pkgName] {
		return nil
	}

	pkg, err := buildContext.Import(pkgName, root, 0)
	if err != nil {
		return fmt.Errorf("failed to import %s: %s", pkgName, err)
	}

	if isIgnored(pkg) {
		return nil
	}

	pkgs[pkg.ImportPath] = pkg

	// Don't worry about dependencies for stdlib packages
	if pkg.Goroot && !*delveGoroot {
		return nil
	}

	for _, imp := range getImports(pkg) {
		if _, ok := pkgs[imp]; !ok {
			if err := processPackage(root, imp, level); err != nil {
				return err
			}
		}
	}
	return nil
}

func getImports(pkg *build.Package) []string {
	allImports := pkg.Imports
	if *includeTests {
		allImports = append(allImports, pkg.TestImports...)
		allImports = append(allImports, pkg.XTestImports...)
	}
	var imports []string
	found := make(map[string]struct{})
	for _, imp := range allImports {
		if imp == pkg.ImportPath {
			// Don't draw a self-reference when foo_test depends on foo.
			continue
		}
		if _, ok := found[imp]; ok {
			continue
		}
		found[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return imports
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

func getUrl(name string) (url string) {
	if !strings.HasPrefix(name, "github.") {
		url = fmt.Sprintf("https://godoc.org/%s", name)
		return
	}
	tokens := strings.Split(name, "/")
	if len(tokens) < 3 {
		return
	}
	url = fmt.Sprintf("https://%s/%s/%s", tokens[0], tokens[1], tokens[2])
	return
}

func hasPrefixes(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func isIgnored(pkg *build.Package) bool {
	if len(onlyPrefixes) > 0 && !hasPrefixes(pkg.ImportPath, onlyPrefixes) {
		return true
	}
	return ignored[pkg.ImportPath] || (pkg.Goroot && *ignoreStdlib) || hasPrefixes(pkg.ImportPath, ignoredPrefixes)
}

func printf(format string, args ...interface{}) {
	fmt.Fprintf(output, format, args...)
}

func println(args ...interface{}) {
	fmt.Fprintln(output, args...)
}
