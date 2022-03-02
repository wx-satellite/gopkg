package ast

import (
	"fmt"
	"go/ast"
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
func (g *generator)addDeferStmt(fd *ast.FuncDecl) {

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

	return
}