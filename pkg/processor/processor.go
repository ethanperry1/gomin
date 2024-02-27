// Package processor processes files, subtracting away only the functions declared in the file.
// These functions can be defined through a function declaration syntax or as an inline definition in a variable declaration.
package processor

import (
	"go/ast"
	"go/token"

	"github.com/ethanperry1/gomin/pkg/declarations"
	"github.com/ethanperry1/gomin/pkg/tokens"
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

func (processor *FileProcessor) processFuncLit(name string, lit *ast.FuncLit) *declarations.Decl {
	pos := processor.fst.Position(lit.Body.Pos())
	return &declarations.Decl{
		Line:   pos.Line,
		Column: pos.Column,
		Name:   name,
	}
}

func (processor *FileProcessor) processValueSpec(spec *ast.ValueSpec) []*declarations.Decl {
	if len(spec.Values) == 0 {
		return []*declarations.Decl{}
	}

	decls := processor.processExprs(spec.Values)

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

func (processor *FileProcessor) processGenDecl(genDecl *ast.GenDecl) []*declarations.Decl {
	var decls []*declarations.Decl
	for _, spec := range genDecl.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			// Assign all values to same genDecl comments, if they do not have their own.
			if s.Doc == nil {
				s.Doc = genDecl.Doc
			}
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

func (processor *FileProcessor) processExprs(exprs []ast.Expr) []*declarations.Decl {
	var decls []*declarations.Decl
	for _, expr := range exprs {
		decls = append(decls, processor.processExpr(expr)...)
	}

	return decls
}

func (processor *FileProcessor) processExpr(expr ast.Expr) []*declarations.Decl {
	switch e := expr.(type) {
	case *ast.FuncLit:
		return processor.parseFuncLit(e)
	case *ast.CompositeLit:
		return processor.parseCompositeLit(e)
	case *ast.ParenExpr:
		return processor.parseParenExpr(e)
	case *ast.IndexExpr:
		return processor.parseIndexExpr(e)
	case *ast.IndexListExpr:
		return processor.parseIndexListExpr(e)
	case *ast.SliceExpr:
		return processor.parseSliceExpr(e)
	case *ast.TypeAssertExpr:
		return processor.parseTypeAssertExpr(e)
	case *ast.CallExpr:
		return processor.parseCallExpr(e)
	case *ast.StarExpr:
		return processor.parseStarExpr(e)
	case *ast.UnaryExpr:
		return processor.parseUnaryExpr(e)
	case *ast.BinaryExpr:
		return processor.parseBinaryExpr(e)
	case *ast.KeyValueExpr:
		return processor.parseKeyValueExpr(e)
	case *ast.Ellipsis:
		return processor.parseEllipsis(e)
	}

	return []*declarations.Decl{}
}

func (processor *FileProcessor) parseFuncLit(expr *ast.FuncLit) []*declarations.Decl {
	pos := processor.fst.Position(expr.Body.Pos())
	return []*declarations.Decl{
		{
			Line:   pos.Line,
			Column: pos.Column,
			Name:   "no name (function literal)",
		},
	}
}

func (processor *FileProcessor) parseCompositeLit(expr *ast.CompositeLit) []*declarations.Decl {
	return processor.processExprs(expr.Elts)
}

func (processor *FileProcessor) parseParenExpr(expr *ast.ParenExpr) []*declarations.Decl {
	return processor.processExpr(expr.X)
}

func (processor *FileProcessor) parseIndexExpr(expr *ast.IndexExpr) []*declarations.Decl {
	return append(processor.processExpr(expr.X), processor.processExpr(expr.Index)...)
}

func (processor *FileProcessor) parseIndexListExpr(expr *ast.IndexListExpr) []*declarations.Decl {
	return append(processor.processExpr(expr.X), processor.processExprs(expr.Indices)...)
}

func (processor *FileProcessor) parseSliceExpr(expr *ast.SliceExpr) []*declarations.Decl {
	decls := processor.processExpr(expr.X)

	if expr.Low != nil {
		decls = append(decls, processor.processExpr(expr.Low)...)
	}

	if expr.High != nil {
		decls = append(decls, processor.processExpr(expr.High)...)
	}

	if expr.Max != nil {
		decls = append(decls, processor.processExpr(expr.Max)...)
	}

	return decls
}

func (processor *FileProcessor) parseTypeAssertExpr(expr *ast.TypeAssertExpr) []*declarations.Decl {
	decls := processor.processExpr(expr.X)
	if expr.Type != nil {
		decls = append(decls, processor.processExpr(expr.Type)...)
	}

	return decls
}

func (processor *FileProcessor) parseCallExpr(expr *ast.CallExpr) []*declarations.Decl {
	return append(processor.processExpr(expr.Fun), processor.processExprs(expr.Args)...)
}

func (processor *FileProcessor) parseStarExpr(expr *ast.StarExpr) []*declarations.Decl {
	return processor.processExpr(expr.X)
}

func (processor *FileProcessor) parseUnaryExpr(expr *ast.UnaryExpr) []*declarations.Decl {
	return processor.processExpr(expr.X)
}

func (processor *FileProcessor) parseBinaryExpr(expr *ast.BinaryExpr) []*declarations.Decl {
	return append(processor.processExpr(expr.X), processor.processExpr(expr.Y)...)
}

func (processor *FileProcessor) parseKeyValueExpr(expr *ast.KeyValueExpr) []*declarations.Decl {
	return append(processor.processExpr(expr.Key), processor.processExpr(expr.Value)...)
}

func (processor *FileProcessor) parseEllipsis(expr *ast.Ellipsis) []*declarations.Decl {
	return processor.processExpr(expr.Elt)
}
