/*
 * Copyright 2018. bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Generator struct {
	buf     bytes.Buffer
	db      *sql.DB
	dialect Dialect
	dir     string
	ctx     *Context
}

func NewGenerator(dir string, db *sql.DB, dialect Dialect) *Generator {
	g := &Generator{
		db:      db,
		dialect: dialect,
		dir:     dir,
		ctx:     NewContext(),
	}
	return g
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

func (g *Generator) Println(a ...interface{}) {
	fmt.Fprintln(&g.buf, a...)
}

func (g *Generator) ParserPkg() {
	var (
		files []*ast.File
		err   error
	)
	fs := token.NewFileSet()

	pkgMap, err := parser.ParseDir(fs, g.dir, nil, 0)
	if err != nil {
		g.ctx.Error(err)
		return
	}
	for _, pkg := range pkgMap {
		for filename, f := range pkg.Files {
			_, filename = path.Split(filename)
			if filename != "toyorm_models.go" {
				files = append(files, f)
			}
		}
	}
	config := types.Config{Importer: importer.For("source", nil), FakeImportC: true}
	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}
	pkg, err := config.Check(g.dir, fs, files, info)

	if err != nil {
		g.ctx.Error(err)
		return
	}
	if pkg.Name() == "" {
		g.ctx.Error(errors.New("nil package name"))
		return
	}
	g.ctx.Pkg = pkg.Name()
	for _, elem := range pkg.Scope().Names() {
		g.ctx.Used[elem] = true
	}
}

func (g *Generator) Generate() {
	g.Printf("// Code generated by \"toy-factory %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Println()
	g.Printf("package %s", g.ctx.Pkg)
	g.Println()
	g.dialect.ParserTables(g.db, nil, g.ctx)
	if g.ctx.error != nil {
		return
	}
	// add import package
	var importList []string
	for imp, alias := range g.ctx.Import {
		if alias == "" {
			importList = append(importList, fmt.Sprintf(`"%s"`, imp))
		} else {
			importList = append(importList, fmt.Sprintf(`%s "%s"`, alias, imp))
		}
	}
	if len(g.ctx.Import) != 0 {
		g.Printf("import (")
		g.Println()
		for imp, alias := range g.ctx.Import {
			if alias == "" {
				g.Printf(`"%s"`, imp)
			} else {
				g.Printf(`%s "%s"`, alias, imp)
			}
			g.Println()
		}
		g.Println(")")
	}

	for _, model := range g.ctx.Models {
		if g.ctx.Used[model.Name] {
			g.ctx.Error(errors.New(fmt.Sprintf("model name %s already exist in globol scope", model.Name)))
			return
		}
		g.Println(model.String())
	}
}

func (g *Generator) Output() {
	data, err := format.Source(g.buf.Bytes())
	if err != nil {
		g.ctx.Error(err)
		return
	}
	filename := "toyorm_models.go"
	err = ioutil.WriteFile(path.Join(g.dir, filename), data, 0644)
	if err != nil {
		g.ctx.Error(err)
		return
	}
}

func (g *Generator) Run() error {

	g.ParserPkg()
	if g.ctx.error != nil {
		return g.ctx.error
	}
	g.Generate()
	if g.ctx.error != nil {
		return g.ctx.error
	}
	g.Output()
	if g.ctx.error != nil {
		return g.ctx.error
	}
	return nil
}