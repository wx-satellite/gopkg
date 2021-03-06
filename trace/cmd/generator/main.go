package main

import (
	"flag"
	"fmt"
	"github.com/wx-satellite/gopkg/trace/generator/ast"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	wrote bool
)

func init() {
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

	generator := ast.New("github.com/wx-satellite/gopkg/trace","trace","Do")

	res,err := generator.Generate(file)
	if err != nil {
		panic(err)
	}

	if len(res) ==0 {
		fmt.Println("添加 trace 失败")
		return
	}

	if !wrote {
		fmt.Println(string(res))
		return
	}


	// 写入源文件
	if err = ioutil.WriteFile(file,res,0666); err != nil {
		fmt.Println("写入文件失败")
		return
	}


	fmt.Println("添加 trace 成功")

}