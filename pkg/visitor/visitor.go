//gobar:min:1.00
package visitor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"gobar/pkg/tokens"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type (
	Parser interface {
		ParseComment(comment string) (tokens.Comparer, error)
	}
	Emplacer interface {
		Emplace(fileName string, file *File)
	}
	Processor interface {
		Process(comments []*ast.Comment) ([]tokens.Comparer, error)
	}
)

type File struct {
	Ast *ast.File
	Fst *token.FileSet
}

type FileEmplacer struct {
	files     map[string][]tokens.Comparer
	funcs     map[token.Position][]tokens.Comparer
	asts 	  map[string]*File
}

func NewEmplacer(asts map[string]*File) *FileEmplacer {
	return &FileEmplacer{
		asts: asts,
	}
}

func (emplacer *FileEmplacer) Emplace(fileName string, file *File) {
	emplacer.asts[fileName] = file
}

type Visitor struct {
	emplacer     Emplacer
	Readdirnames func(n int) (names []string, err error)
}

func NewVisitor(emplacer Emplacer, options ...func(*Visitor)) *Visitor {
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
			fst := token.NewFileSet()
			tree, err := parser.ParseFile(fst, joined, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			visitor.emplacer.Emplace(joined, &File{
				Ast: tree,
				Fst: fst,
			})
		}
	}

	return nil
}
