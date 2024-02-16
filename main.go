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

//gobar:0.9
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

type Emplacer interface {
	Emplace(file string, ast *ast.File)
}

type FileEmplacer struct {
	asts map[string]*ast.File
}

func NewEmplacer(asts map[string]*ast.File) *FileEmplacer {
	return &FileEmplacer{
		asts: asts,
	}
}

func (emplacer *FileEmplacer) Emplace(file string, fileTree *ast.File) {
	emplacer.asts[file] = fileTree
	for _, decl := range fileTree.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		fmt.Printf("%d\n", funcDecl.Body.Pos())
	}
}

type Visitor struct {
	emplacer     Emplacer
	Readdirnames func(n int) (names []string, err error)
}

func NewVisitor(emplacer     Emplacer, options ...func(*Visitor)) *Visitor {
	visitor := &Visitor{
		emplacer: emplacer,
	}

	for _, option := range options {
		option(visitor)
	}

	return visitor
}

func (visitor *Visitor) Visit(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	names, err := file.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		// Ignore nested modules.
		if name == "go.mod" {
			return nil
		}
	}

	for _, name := range names {
		if strings.HasSuffix(name, ".go") {
			joined := filepath.Join(path, name)
			f, err := parser.ParseFile(token.NewFileSet(), joined, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			visitor.emplacer.Emplace(joined, f)
		}
	}

	return nil
}
