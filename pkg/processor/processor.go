// Package processor processes files, subtracting away only the functions declared in the file.
// These functions can be defined through a function declaration syntax or as an inline definition in a variable declaration.
package processor

import (
	"go/ast"
	"go/token"
	"gobar/pkg/declarations"
	"gobar/pkg/tokens"
)

type (
	Parser interface {
		ParseComment(comment string, parents ...tokens.Comparer) (tokens.Comparer, error)
	}
)

type FileProcessor struct {
	fst       *token.FileSet
	tree      *ast.File
	directory string
	fileName  string
}

func New(
	fst *token.FileSet,
	tree *ast.File,
	directory string,
	fileName string,
) *FileProcessor {
	return &FileProcessor{
		fst:       fst,
		tree:      tree,
		directory: directory,
		fileName:  fileName,
	}
}

func (processor *FileProcessor) processFuncTypeSpec(spec *ast.ValueSpec) []*declarations.Decl {

	decls := make([]*declarations.Decl, len(spec.Names))
	for idx, name := range spec.Names {
		decls[idx].Name = name.Name
	}

	for idx, value := range spec.Values {
		pos := processor.fst.Position(value.Pos())
		decls[idx].Line = pos.Line
		decls[idx].Column = pos.Column
	}

	if spec.Doc != nil {
		comments := make([]string, len(spec.Doc.List))
		for idx, comment := range spec.Doc.List {
			comments[idx] = comment.Text
		}

		for _, decl := range decls {
			decl.Comments = comments
		}
	}

	return decls
}

func (processor *FileProcessor) processValueSpec(spec *ast.ValueSpec) []*declarations.Decl {
	switch spec.Type.(type) {
	case *ast.FuncType:
		return processor.processFuncTypeSpec(spec)
	default:
		return []*declarations.Decl{}
	}
}

func (processor *FileProcessor) processGenDecl(genDecl *ast.GenDecl) []*declarations.Decl {
	var decls []*declarations.Decl
	for _, spec := range genDecl.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			decls = append(decls, processor.processValueSpec(s)...)
		}
	}

	return decls
}

func (processor *FileProcessor) processFuncDecl(funcDecl *ast.FuncDecl) []*declarations.Decl {
	pos := processor.fst.Position(funcDecl.Body.Pos())

	decl := &declarations.Decl{
		Line:   pos.Line,
		Column: pos.Column,
		Name:   funcDecl.Name.Name,
	}

	if funcDecl.Doc != nil {
		comments := make([]string, len(funcDecl.Doc.List))
		for idx, comment := range funcDecl.Doc.List {
			comments[idx] = comment.Text
		}
		decl.Comments = comments
	}

	return []*declarations.Decl{decl}
}

func (processor *FileProcessor) Process() []*declarations.Decl {
	var decls []*declarations.Decl
	for _, decl := range processor.tree.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			decls = append(decls, processor.processGenDecl(d)...)
		case *ast.FuncDecl:
			decls = append(decls, processor.processFuncDecl(d)...)
		}
	}

	// Append filename and directory to all declarations.
	for _, decl := range decls {
		decl.FileName = processor.fileName
		decl.Directory = processor.directory
	}

	return decls
}
