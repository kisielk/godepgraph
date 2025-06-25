package main

import (
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type graphPrinter interface {
	writeHeader(hLayout bool)
	writeNode(pkgName string, attrs *build.Package)
	writeEdge(u string, v string)
	writeEnd()
}

var (
	pkgs        map[string]*build.Package
	erroredPkgs map[string]bool

	ignored = map[string]bool{
		"C": true,
	}
	ignoredPrefixes []string
	onlyPrefixes    []string

	ignoreStdlib   = flag.Bool("nostdlib", false, "ignore packages in the Go standard library")
	ignoreVendor   = flag.Bool("novendor", false, "ignore packages in the vendor directory")
	stopOnError    = flag.Bool("stoponerror", true, "stop on package import errors")
	withGoroot     = flag.Bool("withgoroot", false, "show dependencies of packages in the Go standard library")
	ignorePrefixes = flag.String("ignoreprefixes", "", "a comma-separated list of prefixes to ignore")
	ignorePackages = flag.String("ignorepackages", "", "a comma-separated list of packages to ignore")
	onlyPrefix     = flag.String("onlyprefixes", "", "a comma-separated list of prefixes to include")
	tagList        = flag.String("tags", "", "a comma-separated list of build tags to consider satisfied during the build")
	horizontal     = flag.Bool("horizontal", false, "lay out the dependency graph horizontally instead of vertically")
	withTests      = flag.Bool("withtests", false, "include test packages")
	maxLevel       = flag.Int("maxlevel", 256, "max level of go dependency graph")
	format         = flag.String("format", "dot", "output format of graph (dot, mermaid)")

	buildTags    []string
	buildContext = build.Default
)

func init() {
	flag.BoolVar(ignoreStdlib, "s", false, "(alias for -nostdlib) ignore packages in the Go standard library")
	flag.StringVar(ignorePrefixes, "p", "", "(alias for -ignoreprefixes) a comma-separated list of prefixes to ignore")
	flag.StringVar(ignorePackages, "i", "", "(alias for -ignorepackages) a comma-separated list of packages to ignore")
	flag.StringVar(onlyPrefix, "o", "", "(alias for -onlyprefixes) a comma-separated list of prefixes to include")
	flag.BoolVar(withTests, "t", false, "(alias for -withtests) include test packages")
	flag.IntVar(maxLevel, "l", 256, "(alias for -maxlevel) maximum level of the go dependency graph")
	flag.BoolVar(withGoroot, "d", false, "(alias for -withgoroot) show dependencies of packages in the Go standard library")
}

func main() {
	pkgs = make(map[string]*build.Package)
	erroredPkgs = make(map[string]bool)
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
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

	printer, err := getPrinter()
	if err != nil {
		log.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd: %s", err)
	}

	cmdArgs := []string{"list"}
	if *tagList != "" {
		cmdArgs = append(cmdArgs, "-tags", *tagList)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("go", cmdArgs...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	resolved := strings.Fields(string(out))

	for _, pkg := range resolved {
		if err := processPackage(cwd, pkg, 0, "", *stopOnError); err != nil {
			log.Fatal(err)
		}
	}

	printer.writeHeader(*horizontal)

	// sort packages
	pkgKeys := []string{}
	for k := range pkgs {
		pkgKeys = append(pkgKeys, k)
	}
	sort.Strings(pkgKeys)

	for _, pkgName := range pkgKeys {
		pkg := pkgs[pkgName]

		if isIgnored(pkg) {
			continue
		}

		printer.writeNode(pkgName, pkg)

		// Don't render imports from packages in Goroot
		if pkg.Goroot && !*withGoroot {
			continue
		}

		for _, imp := range getImports(pkg) {
			impPkg := pkgs[imp]
			if impPkg == nil || isIgnored(impPkg) {
				continue
			}

			printer.writeEdge(pkgName, imp)
		}
	}

	printer.writeEnd()
}

func pkgDocsURL(pkgName string) string {
	return "https://godoc.org/" + pkgName
}

func processPackage(root string, pkgName string, level int, importedBy string, stopOnError bool) error {
	if level++; level > *maxLevel {
		return nil
	}
	if ignored[pkgName] {
		return nil
	}

	pkg, buildErr := buildContext.Import(pkgName, root, 0)
	if buildErr != nil {
		if stopOnError {
			return fmt.Errorf("failed to import %s (imported at level %d by %s):\n%s", pkgName, level, importedBy, buildErr)
		}
	}

	if isIgnored(pkg) {
		return nil
	}

	importPath := normalizeVendor(pkgName)
	if buildErr != nil {
		erroredPkgs[importPath] = true
	}

	pkgs[importPath] = pkg

	// Don't worry about dependencies for stdlib packages
	if pkg.Goroot && !*withGoroot {
		return nil
	}

	for _, imp := range getImports(pkg) {
		if _, ok := pkgs[imp]; !ok {
			if err := processPackage(pkg.Dir, imp, level, pkgName, stopOnError); err != nil {
				return err
			}
		}
	}
	return nil
}

func getImports(pkg *build.Package) []string {
	allImports := pkg.Imports
	if *withTests {
		allImports = append(allImports, pkg.TestImports...)
		allImports = append(allImports, pkg.XTestImports...)
	}
	var imports []string
	found := make(map[string]struct{})
	for _, imp := range allImports {
		if imp == normalizeVendor(pkg.ImportPath) {
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

func hasPrefixes(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func isIgnored(pkg *build.Package) bool {
	if len(onlyPrefixes) > 0 && !hasPrefixes(normalizeVendor(pkg.ImportPath), onlyPrefixes) {
		return true
	}

	if *ignoreVendor && isVendored(pkg.ImportPath) {
		return true
	}
	return ignored[normalizeVendor(pkg.ImportPath)] || (pkg.Goroot && *ignoreStdlib) || hasPrefixes(normalizeVendor(pkg.ImportPath), ignoredPrefixes)
}

func hasBuildErrors(pkg *build.Package) bool {
	if len(erroredPkgs) == 0 {
		return false
	}

	v, ok := erroredPkgs[normalizeVendor(pkg.ImportPath)]
	if !ok {
		return false
	}
	return v
}

func isVendored(path string) bool {
	return strings.Contains(path, "/vendor/")
}

func normalizeVendor(path string) string {
	pieces := strings.Split(path, "vendor/")
	return pieces[len(pieces)-1]
}

func getPrinter() (graphPrinter, error) {
	switch *format {
	case "dot":
		return newGraphvizPrinter(), nil
	case "mermaid":
		return newMermaidPrinter(), nil
	default:
		return nil, fmt.Errorf("invalid format flag %q, must be: dot, mermaid", *format)
	}
}
