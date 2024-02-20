package evaluate

import (
	"gobar/pkg/coverage"
	"gobar/pkg/declarations"
	"gobar/pkg/parser"
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

	// Create new directive parser, which will check each project file for gobar coverage directives.
	psr := parser.New(tokens.CreatorsByCommand)

	for dir, files := range dirs {

		pack := unit.NewPackage(dir)
		project.WithChild(pack)

		var packageComments []string
		for _, file := range files {
			if file.Ast.Doc != nil {
				for _, comment := range file.Ast.Doc.List {
					packageComments = append(packageComments, comment.Text)
				}
			}
		}

		packDirectives, err := psr.ParsePkgComments(packageComments)
		if err != nil {
			return nil, &InvalidPackageDirectiveError{
				name: dir,
				err:  err,
			}
		}

		pack.WithDirectives(packDirectives...)

		for name, file := range files {

			// Continue if no profile is found.
			profile := profilesByName.Get(filepath.Join(evaluator.name, dir, name))
			if profile == nil {
				continue
			}

			fl := unit.NewFile(name)
			pack.WithChild(fl)

			var fileComments []string
			if file.Ast.Doc != nil {
				for _, comment := range file.Ast.Doc.List {
					fileComments = append(fileComments, comment.Text)
				}
			}

			fileDirectives, err := psr.ParseComments(fileComments, packDirectives...)
			if err != nil {
				return nil, &InvalidFileDirectiveError{
					err:  err,
					name: name,
					dir:  dir,
				}
			}

			fl.WithDirectives(fileDirectives...)

			decls := processor.New(file.Fst, file.Ast, dir, name).Process()

			// Sort all declarations.
			sortedDeclarations := declarations.New(declarations.Sort(decls))

			coverageCalculator := coverage.New(sortedDeclarations)

			results := coverageCalculator.ProcessCoverage(profile)

			for _, decl := range decls {
				cov := results[decl.Name]
				block := unit.NewBlock(decl.Name, cov)
				fl.WithChild(block)

				blockDirectives, err := psr.ParseComments(decl.Comments, fileDirectives...)
				if err != nil {
					return nil, &InvalidBlockDirectiveError{
						err:   err,
						dir:   dir,
						name:  name,
						block: decl.Name,
					}
				}

				block.WithDirectives(blockDirectives...)
			}
		}
	}

	return project.Evaluate()
}
