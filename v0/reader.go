package api

import (
	"os"
	"path/filepath"

	"github.com/ethanperry1/gomin/pkg/declarations"
	"github.com/ethanperry1/gomin/pkg/profiles"
	"github.com/ethanperry1/gomin/pkg/visitor"

	"golang.org/x/tools/cover"
)

type ProfileReader struct {
	profile string
	name    string
	root    string
}

func NewProfileReader(
	root string,
	profile string,
	name string,
) *ProfileReader {
	return &ProfileReader{
		root:    root,
		profile: profile,
		name:    name,
	}
}

func (reader *ProfileReader) CreateNodeTree() (Node, error) {
	err := os.Chdir(reader.root)
	if err != nil {
		return nil, err
	}

	profs, err := cover.ParseProfiles(reader.profile)
	if err != nil {
		return nil, err
	}

	profilesByName := profiles.New(profs)

	dirs := make(map[string]map[string]*visitor.File)
	emplacer := visitor.NewEmplacer(dirs)
	visitor := visitor.NewVisitor(emplacer)

	err = filepath.WalkDir(".", visitor.Visit)
	if err != nil {
		return nil, err
	}

	pkgCount := 0
	pkgOptions := make([]func(*node), len(dirs))
	for dir, files := range dirs {
		var fileOptions []func(*node)
		for name, file := range files {
			profile := profilesByName.Get(filepath.Join(reader.name, dir, name))
			if profile == nil {
				continue
			}

			decls := NewFileProcessor(file.Fst, file.Ast, dir, name).Process()

			sortedDeclarations := declarations.New(declarations.Sort(decls))

			funcStmts := make(map[any]*statements)
			for _, block := range profile.Blocks {
				decl := sortedDeclarations.DeclByPosition(block.StartLine, block.StartCol)
				result, ok := funcStmts[decl]
				if !ok {
					result = &statements{}
					funcStmts[decl] = result
				}

				result.total += block.NumStmt

				if block.Count > 0 {
					result.covered += block.NumStmt
				}
			}

			count := 0
			blockOptions := make([]func(*node), len(funcStmts))
			for k, v := range funcStmts {
				blockOptions[count] = AddNode(k, NewNode(AddStatement(v)))
				count++
			}

			fileOptions = append(fileOptions, AddNode(name, NewNode(blockOptions...)))
		}

		pkgOptions[pkgCount] = AddNode(dir, NewNode(fileOptions...))
		pkgCount++
	}

	return NewNode(pkgOptions...), nil
}
