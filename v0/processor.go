// package v0 processes files, subtracting away only the functions declared in the file.
// These functions can be defined through a function declaration syntax or as an inline definition in a variable declaration.
package v0

import (
	"go/ast"
	"go/token"

	"github.com/ethanperry1/gomin/pkg/declarations"
)

type FileProcessor struct {
	fst       *token.FileSet
	tree      *ast.File
	directory string
	fileName  string
}

func NewFileProcessor(
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

func (processor *FileProcessor) processValueSpec(index int, spec *ast.ValueSpec) (int, []*declarations.Decl) {
	if len(spec.Values) == 0 {
		return index, []*declarations.Decl{}
	}

	index, decls := processor.processExprs(index, spec.Values)

	if spec.Doc != nil {
		comments := make([]string, len(spec.Doc.List))
		for idx, comment := range spec.Doc.List {
			comments[idx] = comment.Text
		}

		for _, decl := range decls {
			decl.Comments = comments
		}
	}

	return index, decls
}

func (processor *FileProcessor) processGenDecl(index int, genDecl *ast.GenDecl) (int, []*declarations.Decl) {
	var allDecls []*declarations.Decl
	var decls []*declarations.Decl
	for _, spec := range genDecl.Specs {
		switch s := spec.(type) {
		case *ast.ValueSpec:
			// Assign all values to same genDecl comments, if they do not have their own.
			if s.Doc == nil {
				s.Doc = genDecl.Doc
			}
			index, decls = processor.processValueSpec(index, s)
			allDecls = append(allDecls, decls...)
		}
	}

	return index, allDecls
}

func (processor *FileProcessor) processFuncDecl(funcDecl *ast.FuncDecl) []*declarations.Decl {
	pos := processor.fst.Position(funcDecl.Body.Pos())

	decl := &declarations.Decl{
		Line:   pos.Line,
		Column: pos.Column,
		Name:   funcDecl.Name.Name,
	}

	if funcDecl.Recv != nil {
		decl.Name = NamePair{
			X: processor.parseField(funcDecl.Recv.List[0]),
			Y: funcDecl.Name.Name,
		}
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
	index := 0
	var allDecls []*declarations.Decl
	var decls []*declarations.Decl
	for _, decl := range processor.tree.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			index, decls = processor.processGenDecl(index, d)
			allDecls = append(allDecls, decls...)
		case *ast.FuncDecl:
			allDecls = append(allDecls, processor.processFuncDecl(d)...)
		}
	}

	// Append filename and directory to all declarations.
	for _, decl := range allDecls {
		decl.FileName = processor.fileName
		decl.Directory = processor.directory
	}

	return allDecls
}

func (processor *FileProcessor) processExprs(index int, exprs []ast.Expr) (int, []*declarations.Decl) {
	var allDecls []*declarations.Decl
	var decls []*declarations.Decl
	for _, expr := range exprs {
		index, decls = processor.processExpr(index, expr)
		allDecls = append(allDecls, decls...)
	}

	return index, allDecls
}

func (processor *FileProcessor) processExpr(index int, expr ast.Expr) (int, []*declarations.Decl) {
	switch e := expr.(type) {
	case *ast.FuncLit:
		return processor.parseFuncLit(index, e)
	case *ast.CompositeLit:
		return processor.parseCompositeLit(index, e)
	case *ast.ParenExpr:
		return processor.parseParenExpr(index, e)
	case *ast.IndexExpr:
		return processor.parseIndexExpr(index, e)
	case *ast.IndexListExpr:
		return processor.parseIndexListExpr(index, e)
	case *ast.SliceExpr:
		return processor.parseSliceExpr(index, e)
	case *ast.TypeAssertExpr:
		return processor.parseTypeAssertExpr(index, e)
	case *ast.CallExpr:
		return processor.parseCallExpr(index, e)
	case *ast.StarExpr:
		return processor.parseStarExpr(index, e)
	case *ast.UnaryExpr:
		return processor.parseUnaryExpr(index, e)
	case *ast.BinaryExpr:
		return processor.parseBinaryExpr(index, e)
	case *ast.KeyValueExpr:
		return processor.parseKeyValueExpr(index, e)
	case *ast.Ellipsis:
		return processor.parseEllipsis(index, e)
	}

	return index, []*declarations.Decl{}
}

func (processor *FileProcessor) parseFuncLit(index int, expr *ast.FuncLit) (int, []*declarations.Decl) {
	pos := processor.fst.Position(expr.Body.Pos())
	return index + 1, []*declarations.Decl{
		{
			Line:   pos.Line,
			Column: pos.Column,
			Name:   index,
		},
	}
}

func (processor *FileProcessor) parseCompositeLit(index int, expr *ast.CompositeLit) (int, []*declarations.Decl) {
	return processor.processExprs(index, expr.Elts)
}

func (processor *FileProcessor) parseParenExpr(index int, expr *ast.ParenExpr) (int, []*declarations.Decl) {
	return processor.processExpr(index, expr.X)
}

func (processor *FileProcessor) parseIndexExpr(index int, expr *ast.IndexExpr) (int, []*declarations.Decl) {
	index, xDecls := processor.processExpr(index, expr.X)
	index, iDecls := processor.processExpr(index, expr.Index)
	return index, append(xDecls, iDecls...)
}

func (processor *FileProcessor) parseIndexListExpr(index int, expr *ast.IndexListExpr) (int, []*declarations.Decl) {
	index, xDecls := processor.processExpr(index, expr.X)
	index, iDecls := processor.processExprs(index, expr.Indices)
	return index, append(xDecls, iDecls...)
}

func (processor *FileProcessor) parseSliceExpr(index int, expr *ast.SliceExpr) (int, []*declarations.Decl) {
	index, decls := processor.processExpr(index, expr.X)
	var extraDecls []*declarations.Decl

	if expr.Low != nil {
		index, extraDecls = processor.processExpr(index, expr.Low)
		decls = append(decls, extraDecls...)
	}

	if expr.High != nil {
		index, extraDecls = processor.processExpr(index, expr.High)
		decls = append(decls, extraDecls...)
	}

	if expr.Max != nil {
		index, extraDecls = processor.processExpr(index, expr.Max)
		decls = append(decls, extraDecls...)
	}

	return index, decls
}

func (processor *FileProcessor) parseTypeAssertExpr(index int, expr *ast.TypeAssertExpr) (int, []*declarations.Decl) {
	index, decls := processor.processExpr(index, expr.X)
	var extraDecls []*declarations.Decl

	if expr.Type != nil {
		index, extraDecls = processor.processExpr(index, expr.Type)
		decls = append(decls, extraDecls...)
	}

	return index, decls
}

func (processor *FileProcessor) parseCallExpr(index int, expr *ast.CallExpr) (int, []*declarations.Decl) {
	index, fDecls := processor.processExpr(index, expr.Fun)
	index, aDecls := processor.processExprs(index, expr.Args)
	return index, append(fDecls, aDecls...)
}

func (processor *FileProcessor) parseStarExpr(index int, expr *ast.StarExpr) (int, []*declarations.Decl) {
	return processor.processExpr(index, expr.X)
}

func (processor *FileProcessor) parseUnaryExpr(index int, expr *ast.UnaryExpr) (int, []*declarations.Decl) {
	return processor.processExpr(index, expr.X)
}

func (processor *FileProcessor) parseBinaryExpr(index int, expr *ast.BinaryExpr) (int, []*declarations.Decl) {
	index, xDecls := processor.processExpr(index, expr.X)
	index, yDecls := processor.processExpr(index, expr.Y)
	return index, append(xDecls, yDecls...)
}

func (processor *FileProcessor) parseKeyValueExpr(index int, expr *ast.KeyValueExpr) (int, []*declarations.Decl) {
	index, kDecls := processor.processExpr(index, expr.Key)
	index, vDecls := processor.processExpr(index, expr.Value)
	return index, append(kDecls, vDecls...)
}

func (processor *FileProcessor) parseEllipsis(index int, expr *ast.Ellipsis) (int, []*declarations.Decl) {
	return processor.processExpr(index, expr.Elt)
}

func (processor *FileProcessor) parseField(field *ast.Field) string {
	return processor.parseFieldExpr(field.Type)
}

func (processor *FileProcessor) parseFieldExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return processor.parseFieldStarExpr(t)
	case *ast.Ident:
		return processor.parseFieldIdent(t)
	case *ast.IndexExpr:
		return processor.parseFieldIndexExpr(t)
	case *ast.IndexListExpr:
		return processor.parseFieldIndexListExpr(t)
	}

	return ""
}

func (processor *FileProcessor) parseFieldStarExpr(expr *ast.StarExpr) string {
	return processor.parseFieldExpr(expr.X)
}

func (processor *FileProcessor) parseFieldIdent(expr *ast.Ident) string {
	return expr.Name
}

func (processor *FileProcessor) parseFieldIndexExpr(expr *ast.IndexExpr) string {
	return processor.parseFieldExpr(expr.X)
}

func (processor *FileProcessor) parseFieldIndexListExpr(expr *ast.IndexListExpr) string {
	return processor.parseFieldExpr(expr.X)
}
