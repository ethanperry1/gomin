package evaluate

import (
	"path/filepath"

	"github.com/ethanperry1/gomin/pkg/coverage"
	"github.com/ethanperry1/gomin/pkg/declarations"
	"github.com/ethanperry1/gomin/pkg/processor"
	"github.com/ethanperry1/gomin/pkg/profiles"
	"github.com/ethanperry1/gomin/pkg/tokens"
	"github.com/ethanperry1/gomin/pkg/unit"
	"github.com/ethanperry1/gomin/pkg/visitor"

	"golang.org/x/tools/cover"
)

type EvaluatorOptions func(*Evaluator)

type Evaluator struct {
	name                  string
	root                  string
	profile               string
	defaultPackageMinimum float64
	defaultFileMinimum    float64
	defaultBlockMinimum   float64
}

func InitDefaultPackageMinimum(defaultMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultPackageMinimum = defaultMinimum
	}
}

func InitDefaultFileMinimum(defaultFileMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultFileMinimum = defaultFileMinimum
	}
}

func InitDefaultFunctionMinimum(defaultFunctionMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultBlockMinimum = defaultFunctionMinimum
	}
}

func New(name, root, profile string, options ...EvaluatorOptions) *Evaluator {
	evalutor := &Evaluator{
		defaultPackageMinimum: 0,
		defaultFileMinimum:    0,
		defaultBlockMinimum:   0,
		name:                  name,
		root:                  root,
		profile:               profile,
	}

	for _, option := range options {
		option(evalutor)
	}

	return evalutor
}

func (evaluator *Evaluator) EvalCoverage() (unit.Coverage, error) {
	// Parse profiles from go cover tool produced coverage report.
	profs, err := cover.ParseProfiles(evaluator.profile)
	if err != nil {
		return nil, err
	}

	profilesByName := profiles.New(profs)

	// Create file visitor to parse abstract syntax tree from each applicable file.
	dirs := make(map[string]map[string]*visitor.File)
	emplacer := visitor.NewEmplacer(dirs)
	visitor := visitor.NewVisitor(emplacer)

	// Walk over each file in the project.
	err = filepath.WalkDir(evaluator.root, visitor.Visit)
	if err != nil {
		return nil, err
	}

	project := unit.NewParent(evaluator.name)

	projectCmps := make(map[tokens.Level][]tokens.LeveledComparer)
	projectCmps[tokens.Package] = append(projectCmps[tokens.Package], &tokens.ComparerWithLevel{
		Comparer: tokens.NewMinimum(evaluator.defaultPackageMinimum),
	})
	projectCmps[tokens.File] = append(projectCmps[tokens.File], &tokens.ComparerWithLevel{
		Comparer: tokens.NewMinimum(evaluator.defaultFileMinimum),
	})
	projectCmps[tokens.Block] = append(projectCmps[tokens.Block], &tokens.ComparerWithLevel{
		Comparer: tokens.NewMinimum(evaluator.defaultBlockMinimum),
	})

	for dir, files := range dirs {

		pkgCmps := make(map[tokens.Level][]tokens.LeveledComparer)

		pack := unit.NewParent(dir)
		project.WithChild(pack)

		var packageTokens [][]string
		for _, file := range files {
			if file.Ast.Doc != nil {
				for _, comment := range file.Ast.Doc.List {
					tks, ok := tokens.Tokenize(comment.Text)
					if ok {
						packageTokens = append(packageTokens, tks)
					}
				}
			}
		}

		for _, tks := range packageTokens {
			cmp, err := tokens.Parse(true, tokens.Package, tks...)
			if err != nil {
				return nil, &InvalidPackageDirectiveError{
					name: dir,
					err:  err,
				}
			}
			if cmp != nil {
				pkgCmps[cmp.Level()] = append(pkgCmps[cmp.Level()], cmp)
			}
		}

		pack.WithDirectives(append(projectCmps[tokens.Package], pkgCmps[tokens.Package]...)...)

		for name, file := range files {

			fileCmps := make(map[tokens.Level][]tokens.LeveledComparer)

			profile := profilesByName.Get(filepath.Join(evaluator.name, dir, name))
			if profile == nil {
				continue
			}

			fl := unit.NewParent(name)
			pack.WithChild(fl)

			var fileTokens [][]string
			if file.Ast.Doc != nil {
				for _, comment := range file.Ast.Doc.List {
					tks, ok := tokens.Tokenize(comment.Text)
					if ok {
						fileTokens = append(fileTokens, tks)
					}
				}
			}

			for _, tks := range fileTokens {
				cmp, err := tokens.Parse(false, tokens.File, tks...)
				if err != nil {
					return nil, &InvalidFileDirectiveError{
						err:  err,
						name: name,
						dir:  dir,
					}
				}

				if cmp != nil {
					fileCmps[cmp.Level()] = append(fileCmps[cmp.Level()], cmp)
				}
			}

			fl.WithDirectives(append(projectCmps[tokens.File], append(pkgCmps[tokens.File], fileCmps[tokens.File]...)...)...)

			decls := processor.New(file.Fst, file.Ast, dir, name).Process()

			sortedDeclarations := declarations.New(declarations.Sort(decls))

			coverageCalculator := coverage.New(sortedDeclarations)

			results := coverageCalculator.ProcessCoverage(profile)

			for _, decl := range decls {

				declCmps := make(map[tokens.Level][]tokens.LeveledComparer)

				cov := results[decl.Name]
				block := unit.NewChild(decl.Name, decl.Line, decl.Column, cov)
				fl.WithChild(block)

				var declTokens [][]string
				for _, comment := range decl.Comments {
					tks, ok := tokens.Tokenize(comment)
					if ok {
						declTokens = append(declTokens, tks)
					}
				}

				for _, tks := range declTokens {
					cmp, err := tokens.Parse(false, tokens.Block, tks...)
					if err != nil {
						return nil, &InvalidBlockDirectiveError{
							err:   err,
							dir:   dir,
							name:  name,
							block: decl.Name,
						}
					}

					if cmp != nil {
						declCmps[cmp.Level()] = append(declCmps[cmp.Level()], cmp)
					}
				}

				block.WithDirectives(append(projectCmps[tokens.Block], append(pkgCmps[tokens.Block], append(fileCmps[tokens.Block], declCmps[tokens.Block]...)...)...)...)
			}
		}
	}

	return project.Evaluate()
}
