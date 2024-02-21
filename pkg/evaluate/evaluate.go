package evaluate

import (
	"gobar/pkg/coverage"
	"gobar/pkg/declarations"
	"gobar/pkg/processor"
	"gobar/pkg/profiles"
	"gobar/pkg/tokens"
	"gobar/pkg/unit"
	"gobar/pkg/visitor"
	"path/filepath"

	"golang.org/x/tools/cover"
)

type EvaluatorOptions func(*Evaluator)

type Evaluator struct {
	name                   string
	root                   string
	profile                string
	defaultMinimum         float64
	defaultFileMinimum     float64
	defaultFunctionMinimum float64
}

func InitDefaultMinimum(defaultMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultMinimum = defaultMinimum
	}
}

func InitDefaultFileMinimum(defaultFileMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultFileMinimum = defaultFileMinimum
	}
}

func InitDefaultFunctionMinimum(defaultFunctionMinimum float64) EvaluatorOptions {
	return func(e *Evaluator) {
		e.defaultFunctionMinimum = defaultFunctionMinimum
	}
}

func New(name, root, profile string, options ...EvaluatorOptions) *Evaluator {
	evalutor := &Evaluator{
		defaultMinimum:         80,
		defaultFileMinimum:     60,
		defaultFunctionMinimum: 40,
		name:                   name,
		root:                   root,
		profile:                profile,
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

	project := unit.NewProject("")

	pkgCmps := make(map[tokens.Level][]tokens.Comparer)

	for dir, files := range dirs {

		pack := unit.NewPackage(dir)
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
			pkgCmps[cmp.Level()] = append(pkgCmps[cmp.Level()], cmp)
		}

		pack.WithDirectives(pkgCmps[tokens.Package]...)

		for name, file := range files {

			fileCmps := make(map[tokens.Level][]tokens.Comparer)

			profile := profilesByName.Get(filepath.Join(evaluator.name, dir, name))
			if profile == nil {
				continue
			}

			fl := unit.NewFile(name)
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
				fileCmps[cmp.Level()] = append(fileCmps[cmp.Level()], cmp)
			}

			fl.WithDirectives(append(pkgCmps[tokens.File], fileCmps[tokens.File]...)...)

			decls := processor.New(file.Fst, file.Ast, dir, name).Process()

			sortedDeclarations := declarations.New(declarations.Sort(decls))

			coverageCalculator := coverage.New(sortedDeclarations)

			results := coverageCalculator.ProcessCoverage(profile)

			for _, decl := range decls {

				declCmps := make(map[tokens.Level][]tokens.Comparer)

				cov := results[decl.Name]
				block := unit.NewBlock(decl.Name, cov)
				fl.WithChild(block)

				var declTokens [][]string
				for _, comment := range decl.Comments {
					tks, ok := tokens.Tokenize(comment)
					if ok {
						declTokens = append(declTokens, tks)
					}
				}

				for _, tks := range declTokens {
					cmp, err := tokens.Parse(false, tokens.File, tks...)
					if err != nil {
						return nil, &InvalidBlockDirectiveError{
							err:   err,
							dir:   dir,
							name:  name,
							block: decl.Name,
						}
					}
					declCmps[cmp.Level()] = append(declCmps[cmp.Level()], cmp)
				}

				block.WithDirectives(append(pkgCmps[tokens.File], append(fileCmps[tokens.File], declCmps[tokens.File]...)...)...)
			}
		}
	}

	return project.Evaluate()
}
