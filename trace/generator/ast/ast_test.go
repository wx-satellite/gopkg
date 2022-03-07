package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

var src = `

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
 ha "golang.org/x/tools/go/ast/astutil"
)

var (
	wrote bool
)

func init() {
	defer ha.usage()
	flag.BoolVar(&wrote,"w",false,"是否将结果写入源文件，默认是控制台")
}

func usage() {
	fmt.Println("执行：instrument [-w] xxx.go 会自动追加trace函数到源文件中")
	flag.PrintDefaults()
}

func main() {
	fmt.Println(os.Args)
	flag.Usage = usage
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("参数不合法")
		usage()
		return
	}

	var file string

	if len(os.Args) == 3{
		file = os.Args[2]
	}

	if len(os.Args) == 2 {
		file = os.Args[1]
	}

	if filepath.Ext(file) !=".go" {
		fmt.Println("仅支持go文件")
		usage()
		return
	}


}
`
func TestFuncDel(t *testing.T) {
	fSet := token.NewFileSet()
	f,_ := parser.ParseFile(fSet, "", src,parser.ParseComments)
	ast.Print(fSet,f)
}


func TestNew(t *testing.T) {
	generator := New("github.com/wx-satellite/gopkg/trace","trace","Do")

	res,err := generator.Generate("./helper.go")

	fmt.Println(string(res))
	fmt.Println(err)
}