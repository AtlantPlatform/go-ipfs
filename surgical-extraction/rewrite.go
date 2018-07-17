// Copyright 2017, 2018 Tensigma Ltd. All rights reserved.
// Use of this source code is governed by Microsoft Reference Source
// License (MS-RSL) that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"
)

// rewriteFile rewrites import statments in the named file
// according to the rules supplied by the map of strings.
//
// Author: https://gist.github.com/jackspirou/61ce33574e9f411b8b4a
func rewriteFile(name string, includes []string, rewriteFn func(string) (string, bool)) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
	if err != nil {
		e := err.Error()
		msg := "expected 'package', found 'EOF'"
		if e[len(e)-len(msg):] == msg {
			return nil
		}
		return err
	}
	change := false
	for _, i := range f.Imports {
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			return err
		}
		if !containsPrefix(path, includes) {
			continue
		}
		if path, ok := rewriteFn(path); ok {
			i.Path.Value = strconv.Quote(path)
			change = true
		}
	}
	for _, cg := range f.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// import \"") {
				ctext := c.Text
				ctext = strings.TrimPrefix(ctext, "// import")
				ctext = strings.TrimSpace(ctext)
				ctext, err := strconv.Unquote(ctext)
				if err != nil {
					return err
				}
				if !containsPrefix(ctext, includes) {
					continue
				}
				if ctext, ok := rewriteFn(ctext); ok {
					c.Text = "// import " + strconv.Quote(ctext)
					change = true
				}
			}
		}
	}
	if !change {
		return nil
	}

	ast.SortImports(fset, f)
	temp := name + ".tmp"
	w, err := os.Create(temp)
	if err != nil {
		return err
	}
	fmtCfg := &printer.Config{
		Tabwidth: 8,
		Mode:     printer.TabIndent | printer.UseSpaces,
	}
	if err := fmtCfg.Fprint(w, fset, f); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return os.Rename(temp, name)
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
