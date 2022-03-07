package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

type generator struct {
	importPath string
	packName string
	funcName string
}

func New(importPath,packName,funcName string) *generator {
	return &generator{
		importPath: importPath,
		packName:   packName,
		funcName:   funcName,
	}
}

func hasFuncDel(f *ast.File) bool {

	for _, decl := range f.Decls {
		if _,ok := decl.(*ast.FuncDecl);ok {
			return true
		}
	}
	return false
}

func (g *generator) addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		// 如果不是函数声明就直接跳过
		if !ok {
			continue
		}
		g.addDeferStmt(fd)
	}
}

// addDeferStmt 向函数声明中追加defer语句
func (g *generator)addDeferStmt(fd *ast.FuncDecl)bool {
	stmts := fd.Body.List
	for _, stmt := range stmts {
		ds,ok:= stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}
		s,ok := ds.Call.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		x,ok := s.X.(*ast.Ident)
		if !ok {
			continue
		}

		// 已经存在
		if x.Name == g.packName && s.Sel.Name == g.funcName {
			return false
		}
	}

	// 不存在，则插入
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{Fun: &ast.SelectorExpr{
			X: &ast.Ident{Name: g.packName},
			Sel: &ast.Ident{Name: g.funcName},
		}},
	}
	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds
	fd.Body.List = newList
	return true
}

func (g *generator) Generate(filename string)(bs []byte,err error) {
	fSet := token.NewFileSet()
	currentAst,err := parser.ParseFile(fSet, filename,nil,parser.ParseComments)
	if err != nil {
		err = fmt.Errorf("error parsing %s: %w", filename, err)
		return
	}

	// 判断是否有函数声明
	if !hasFuncDel(currentAst) {
		return
	}
	// 增加包导入
	astutil.AddImport(fSet,currentAst,g.importPath)

	// 增加 defer trace.Do()()
	g.addDeferTraceIntoFuncDecls(currentAst)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fSet,currentAst)
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}