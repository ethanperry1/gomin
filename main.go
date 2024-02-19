
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

//gobar:min:1.0
func run() error {
	root := os.Getenv("ROOT")
	if root == "" {
		root = "./"
	}

	emplacer := New

	visitor := NewVisitor(make(map[string]*ast.File))

	filepath.WalkDir(root, walker.visit)

	// fst := token.NewFileSet()
	// file := os.Getenv("GOFILE")
	// fmt.Println(file)
	// f, err := parser.ParseFile(fst, file, nil, parser.ParseComments)
	// if err != nil {
	// 	return err
	// }

	// for _, commentGroup := range f.Comments {
	// 	for _, comment := range commentGroup.List {
	// 		fmt.Println(comment.Text)
	// 	}
	// }

	// coverageProf := os.Getenv("PROFILE")
	// fmt.Println(coverageProf)
	// if coverageProf == "" {

	// }

	// profiles, err := cover.ParseProfiles(coverageProf)
	// if err != nil {
	// 	return err
	// }

	// for _, profile := range profiles {
	// 	fmt.Println(profile.FileName)
	// 	for _, block := range profile.Blocks {
	// 		fmt.Printf("%d\n", block.NumStmt)
	// 	}
	// }

	return nil
}

