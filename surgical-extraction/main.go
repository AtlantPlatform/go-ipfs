package main

import (
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/xlab/closer"
)

var app = cli.App("surgical-extraction", "Separates go-ipfs source packages from its repo structure and development flow. Makes IPFS go-gettable.")
var goPaths []string

var (
	targetPkg  = app.StringOpt("P pkg", "github.com/ipfs/go-ipfs/cmd/ipfswatch", "Package of a reference IPFS executable that imports all required IPFS packages.")
	projectOut = app.StringOpt("O out", "bitbucket.org/atlantproject/go-ipfs", "Target must be a project package.")
	stateDB    = app.StringOpt("S state", "extraction.db", "State file (BadgerDB) for relocation bookkeeping.")
	debug      = app.BoolOpt("d debug", true, "Enable debug messages.")
)

func init() {
	log.SetFlags(log.Lshortfile)

	if paths := os.Getenv("GOPATH"); len(paths) == 0 {
		log.Fatalln("GOPATH env variable is not set")
	} else {
		goPaths = strings.Split(paths, ":")
	}
}

func main() {
	app.Command("extract", "Copies sources and rewrites import paths.", extractCmd)
	app.Command("relocate", "Renames a vendored dependency and rewrites the paths.", relocateCmd)
	app.Command("unvendor", "Moves a package from vendored prefix into the project root, rewrites the paths.", unvendorCmd)
	if err := app.Run(os.Args); err != nil {
		log.Fatalln("[ERR]", err)
	}
}

type SourceMap struct {
	OriginSource    SourceFile  `json:"Origin"`
	VendoredSource  *SourceFile `json:"Vendored"`
	RelocatedSource *SourceFile `json:"Relocated"`
}

const (
	SourcesOrigin    string = "Origin"
	SourcesVendored  string = "Vendored"
	SourcesRelocated string = "Relocated"
)

type SourceFile struct {
	FullPath string `json:"Path"`
	Relative string `json:"Relative"`
	Package  string `json:"PackageName"`
}

func (s SourceFile) VendorAuto() SourceFile {
	relativePath := s.Relative
	if s.IsGx() {
		parts := strings.Split(relativePath, "/")
		relativePath = "unknown/" + strings.Join(parts[3:], "/")
	} else {
		relativePath = strings.TrimPrefix(relativePath, "github.com/ipfs/go-ipfs/Godeps/_workspace/src/")
	}
	return SourceFile{
		FullPath: "", // must be composed by writer
		Relative: relativePath,
		Package:  filepath.Dir(relativePath),
	}
}

func (s SourceFile) RelocateAuto(project string) SourceFile {
	if s.IsGodeps() || s.IsGx() {
		return s
	}
	relativePath := strings.TrimPrefix(s.Relative, "github.com/ipfs/go-ipfs/")
	relativePath = fmt.Sprintf("%s/%s", project, relativePath)
	return SourceFile{
		FullPath: "", // must be composed by writer
		Relative: relativePath,
		Package:  filepath.Dir(relativePath),
	}
}

func (s SourceFile) IsGx() bool {
	return strings.HasPrefix(s.Package, "gx/ipfs/")
}

func (s SourceFile) GxRootPackage() (string, bool) {
	if !s.IsGx() {
		return "", false
	}
	parts := strings.Split(s.Package, "/")
	// 0: gx
	// 1: ipfs
	// 2: hash
	// 3: root package
	return parts[3], true
}

func (s SourceFile) IpfsRootPackage() (string, bool) {
	if s.IsGx() || s.IsGodeps() {
		return "", false
	}
	parts := strings.Split(s.Package, "/")
	// 0: github.com
	// 1: ipfs
	// 2: go-ipfs
	// 3: root package
	return "github.com/ipfs/go-ipfs/" + parts[3], true
}

func (s SourceFile) IsGodeps() bool {
	return strings.HasPrefix(s.Package, "github.com/ipfs/go-ipfs/Godeps/_workspace/src/")
}

func (s SourceFile) GodepsRootPackage() (string, bool) {
	if !s.IsGodeps() {
		return "", false
	}
	cleanPrefix := strings.TrimPrefix(s.Package, "github.com/ipfs/go-ipfs/Godeps/_workspace/")
	parts := strings.Split(cleanPrefix, "/")
	// 0: src
	// 1: root repo
	// 2: root author
	// 3: root package
	rootPkg := fmt.Sprintf("%s/%s/%s", parts[1], parts[2], parts[3])
	return rootPkg, true
}

type SourceFiles []SourceFile

func (s SourceFiles) Len() int           { return len(s) }
func (s SourceFiles) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SourceFiles) Less(i, j int) bool { return s[i].FullPath < s[j].FullPath }

func extractCmd(c *cli.Cmd) {
	if len(*projectOut) == 0 {
		closer.Fatalln("[ERR] output project package not specified")
	}
	defaultIncludes := []string{"github.com/ipfs/go-ipfs", "gx/ipfs/"}
	exclude := c.StringsOpt("E exclude", nil, "Exclude specfic package prefixes.")
	include := c.StringsOpt("I include", defaultIncludes, "Include specific package prefixes.")
	reportPath := c.StringOpt("R report", "extraction.json", "Path to extraction report.")
	apply := c.BoolOpt("apply", false, "Copy files over, rewriting targets. Study the extraction map (in .json) before using that.")
	c.Action = func() {
		var projectRoot string
		if pkg := findPackage(*projectOut, *debug); pkg == nil {
			closer.Fatalln("[ERR] failed to find output project package in GOPATH")
		} else {
			projectRoot = pkg.Dir
		}

		defer closer.Close()
		pkg := findPackage(*targetPkg, *debug)
		if pkg == nil {
			closer.Fatalln("[ERR] could not find package in GOPATH:", *targetPkg)
		} else if !pkg.IsCommand() {
			closer.Fatalf("[ERR] package %s is not a command", *targetPkg)
		}
		set, err := findDeps(pkg, *include, *exclude, *debug)
		if err != nil {
			closer.Fatalln("[ERR] failed to read the target package deps:", err)
		}
		deps := make([]SourceFile, 0, len(set))
		for path := range set {
			relative := stripGoPaths(path, goPaths)
			src := SourceFile{
				FullPath: path,
				Relative: relative,
				Package:  filepath.Dir(relative),
			}
			deps = append(deps, src)
		}
		sort.Sort(SourceFiles(deps))

		gxPackages := make(map[string][]SourceFile)
		for _, d := range deps {
			if !d.IsGx() {
				continue
			}
			if rootPkg, ok := d.GxRootPackage(); ok {
				gxPackages[rootPkg] = append(gxPackages[rootPkg], d)
			}
		}
		if *debug {
			log.Println("[INFO] found", len(gxPackages), "deps vendored using Gx")
		}
		godepsPackages := make(map[string][]SourceFile)
		for _, d := range deps {
			if !d.IsGodeps() {
				continue
			}
			if rootPkg, ok := d.GodepsRootPackage(); ok {
				godepsPackages[rootPkg] = append(godepsPackages[rootPkg], d)
			}
		}
		if *debug {
			log.Println("[INFO] found", len(godepsPackages), "deps vendored using Godeps")
		}
		gxPackagesVendored := make(map[string][]SourceFile)
		for pkg, sources := range gxPackages {
			for _, src := range sources {
				src = src.VendorAuto()
				src.FullPath = filepath.Join(projectRoot, "vendor", src.Relative)
				gxPackagesVendored[pkg] = append(gxPackagesVendored[pkg], src)
			}
		}
		godepsPackagesVendored := make(map[string][]SourceFile)
		for pkg, sources := range godepsPackages {
			for _, src := range sources {
				src = src.VendorAuto()
				src.FullPath = filepath.Join(projectRoot, "vendor", src.Relative)
				godepsPackagesVendored[pkg] = append(godepsPackagesVendored[pkg], src)
			}
		}

		ipfsPackages := make(map[string][]SourceFile)
		for _, d := range deps {
			if d.IsGx() || d.IsGodeps() {
				continue
			}
			if rootPkg, ok := d.IpfsRootPackage(); ok {
				ipfsPackages[rootPkg] = append(ipfsPackages[rootPkg], d)
			}
		}
		if *debug {
			log.Println("[INFO] found", len(ipfsPackages), "packages from IPFS")
		}
		ipfsPackagesRelocated := make(map[string][]SourceFile)
		for pkg, sources := range ipfsPackages {
			for _, src := range sources {
				src = src.RelocateAuto(*projectOut)
				cleanPrefix := strings.TrimPrefix(src.Relative, *projectOut+"/")
				src.FullPath = filepath.Join(projectRoot, cleanPrefix)
				ipfsPackagesRelocated[pkg] = append(ipfsPackagesRelocated[pkg], src)
			}
		}

		report := ExtractReport{
			GxResults:     make(map[string]map[string][]SourceFile),
			GodepsResults: make(map[string]map[string][]SourceFile),
			IpfsResults:   make(map[string]map[string][]SourceFile),
		}
		for pkg, sources := range gxPackages {
			report.GxResults[pkg] = map[string][]SourceFile{
				SourcesOrigin:   sources,
				SourcesVendored: gxPackagesVendored[pkg],
			}
		}
		for pkg, sources := range godepsPackages {
			report.GodepsResults[pkg] = map[string][]SourceFile{
				SourcesOrigin:   sources,
				SourcesVendored: godepsPackagesVendored[pkg],
			}
		}
		for pkg, sources := range ipfsPackages {
			report.IpfsResults[pkg] = map[string][]SourceFile{
				SourcesOrigin:    sources,
				SourcesRelocated: ipfsPackagesRelocated[pkg],
			}
		}
		v, _ := json.Marshal(report)
		if err := ioutil.WriteFile(*reportPath, v, 0644); err != nil {
			log.Println("[INFO] failed to write extraction report:", err)
		}
		if !(*apply) {
			return // exit with no further actions
		} else if *debug {
			log.Println("[INFO] copying files over, rewriting paths")
		}
		makeExtract := func(dstPath, srcPath string) {
			if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
				closer.Fatalln("[ERR] failed to ensure a directory:", err)
			}
			if err := copyFile(dstPath, srcPath); err != nil {
				closer.Fatalln("[ERR] failed to copy a file:", err)
			}
			if err := rewriteFile(dstPath, report.Rewrite); err != nil {
				closer.Fatalln("[ERR] failed to rewrite paths in a file:", err)
			}
		}
		var (
			totalPackages int
			totalSources  int
		)
		for _, m := range report.GxResults {
			totalPackages++
			sources := m[SourcesOrigin]
			for i, src := range sources {
				srcPath := src.FullPath
				dstPath := m[SourcesVendored][i].FullPath
				makeExtract(dstPath, srcPath)
				totalSources++
			}
		}
		for _, m := range report.GodepsResults {
			totalPackages++
			sources := m[SourcesOrigin]
			for i, src := range sources {
				srcPath := src.FullPath
				dstPath := m[SourcesVendored][i].FullPath
				makeExtract(dstPath, srcPath)
				totalSources++
			}
		}
		for _, m := range report.IpfsResults {
			totalPackages++
			sources := m[SourcesOrigin]
			for i, src := range sources {
				srcPath := src.FullPath
				dstPath := m[SourcesRelocated][i].FullPath
				makeExtract(dstPath, srcPath)
				totalSources++
			}
		}
		if *debug {
			log.Println("[INFO] relocated", totalPackages, "Go packages")
			log.Println("[INFO] relocated", totalSources, "source files")
		}
	}
}

type ExtractReport struct {
	GxResults     map[string]map[string][]SourceFile `json:"GxResults"`
	GodepsResults map[string]map[string][]SourceFile `json:"GodepsResults"`
	IpfsResults   map[string]map[string][]SourceFile `json:"IpfsResults"`
}

func (e *ExtractReport) Rewrite(path string) (string, bool) {
	if filepath.HasPrefix(path, "gx/ipfs/") {
		for basePkg, m := range e.GxResults {
			sources := m[SourcesOrigin]
			for i, src := range sources {
				if src.Package == path {
					return e.GxResults[basePkg][SourcesVendored][i].Package, true
				}
			}
		}
		log.Println("[WARN] no vendored path for", path)
		return path, false
	}
	if filepath.HasPrefix(path, "github.com/ipfs/go-ipfs/Godeps/_workspace/src/") {
		for basePkg, m := range e.GodepsResults {
			sources := m[SourcesOrigin]
			for i, src := range sources {
				if src.Package == path {
					return e.GodepsResults[basePkg][SourcesVendored][i].Package, true
				}
			}
		}
		log.Println("[WARN] no vendored path for", path)
		return path, false
	}
	for basePkg, m := range e.IpfsResults {
		sources := m[SourcesOrigin]
		for i, src := range sources {
			if src.Package == path {
				return e.GodepsResults[basePkg][SourcesRelocated][i].Package, true
			}
		}
	}
	log.Println("[WARN] no relocation path for", path)
	return path, false
}

func relocateCmd(c *cli.Cmd) {
}

func unvendorCmd(c *cli.Cmd) {
}

func stripGoPaths(path string, gopaths []string) string {
	for _, g := range gopaths {
		path = strings.TrimPrefix(path, fmt.Sprintf("%s/src/", g))
	}
	return path
}

func findPackage(name string, debug bool) *build.Package {
	pkg, err := build.Import(name, "", 0)
	if err != nil {
		if debug {
			log.Println("findPackage:", err)
		}
		return nil
	}
	return pkg
}

func findDeps(p *build.Package, includes, excludes []string, debug bool) (map[string]struct{}, error) {
	excludeRxs, err := compileRxs(excludes)
	if err != nil {
		return nil, err
	}
	deps := make(map[string]struct{}, len(p.Imports))
	addDeps := func(dep *build.Package) {
		if !containsPrefix(dep.ImportPath, includes) {
			return
		}
		for _, f := range dep.GoFiles {
			deps[filepath.Join(dep.Dir, f)] = struct{}{}
		}
		for _, f := range dep.HFiles {
			deps[filepath.Join(dep.Dir, f)] = struct{}{}
		}
		for _, f := range dep.CFiles {
			deps[filepath.Join(dep.Dir, f)] = struct{}{}
		}
	}
	seenDeps := make(map[string]struct{}, len(p.Imports))
	addImports := func(p *build.Package) []*build.Package {
		addDeps(p)
		seenDeps[p.ImportPath] = struct{}{}
		list := make([]*build.Package, 0, len(p.Imports))
		for _, pkg := range p.Imports {
			if _, ok := seenDeps[pkg]; ok {
				continue
			}
			if dep, err := build.Import(pkg, "", 0); err == nil {
				addDeps(dep)
				seenDeps[dep.ImportPath] = struct{}{}
				list = append(list, dep)
			}
		}
		return list
	}
	list := addImports(p)
	for len(list) > 0 {
		newList := make([]*build.Package, 0, len(list))
		for _, p := range list {
			newList = append(newList, addImports(p)...)
		}
		list = newList
	}
	for path := range deps {
		if isMatching(path, excludeRxs) {
			if debug {
				log.Println("gody watch: skipping", path)
			}
			delete(deps, path)
		}
	}
	return deps, nil
}

func isMatching(path string, rxs []*regexp.Regexp) bool {
	for _, rx := range rxs {
		if rx.MatchString(path) {
			return true
		}
	}
	return false
}

func compileRxs(rxs []string) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, rx := range rxs {
		r, err := regexp.Compile(rx)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Regexp: %s error: %v", rx, err)
		}
		compiled = append(compiled, r)
	}
	return compiled, nil
}

func containsPrefix(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if p == "..." {
			return true
		}
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
