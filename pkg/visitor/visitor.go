package visitor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethanperry1/gomin/pkg/tokens"
)

type (
	Parser interface {
		ParseComment(comment string) (tokens.Comparer, error)
	}
	Emplacer interface {
		Emplace(directory string, fileName string, file *File)
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
	asts map[string]map[string]*File
}

func NewEmplacer(asts map[string]map[string]*File) *FileEmplacer {
	return &FileEmplacer{
		asts: asts,
	}
}

func (emplacer *FileEmplacer) Emplace(directory string, fileName string, file *File) {
	dir, ok := emplacer.asts[directory]
	if !ok {
		dir = make(map[string]*File)
		emplacer.asts[directory] = dir
	}

	dir[fileName] = file
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

	// Ignore individual files.
	if !d.IsDir() {
		return nil
	}

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
		// Ignore test files.
		if strings.HasSuffix(name, "_test.go") {
			continue
		}

		if strings.HasSuffix(name, ".go") {
			fst := token.NewFileSet()
			tree, err := parser.ParseFile(fst, filepath.Join(path, name), nil, parser.ParseComments)
			if err != nil {
				return err
			}

			visitor.emplacer.Emplace(path, name, &File{
				Ast: tree,
				Fst: fst,
			})
		}
	}

	return nil
}
