package evaluate

import (
	"gobar/pkg/coverage"
	"gobar/pkg/declarations"
	"gobar/pkg/parser"
	"gobar/pkg/processor"
	"gobar/pkg/profiles"
	"gobar/pkg/tokens"
	"gobar/pkg/visitor"
	"path/filepath"

	"golang.org/x/tools/cover"
)

type EvaluatorOptions func(*Evaluator)

type Evaluator struct {
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

func New(root, profile string, options ...EvaluatorOptions) *Evaluator {
	evalutor := &Evaluator{
		defaultMinimum:         80,
		defaultFileMinimum:     60,
		defaultFunctionMinimum: 40,
		root:                   root,
		profile:                profile,
	}

	for _, option := range options {
		option(evalutor)
	}

	return evalutor
}

func (evaluator *Evaluator) EvalCoverage() error {
	profs, err := cover.ParseProfiles(evaluator.profile)
	if err != nil {
		return err
	}

	profilesByName := profiles.New(profs)

	files := make(map[string]*visitor.File)
	emplacer := visitor.NewEmplacer(files)
	visitor := visitor.NewVisitor(emplacer)

	err = filepath.WalkDir(evaluator.root, visitor.Visit)
	if err != nil {
		return err
	}

	psr := parser.New(tokens.CreatorsByCommand)
	commentProcessor := parser.NewCommentParser(psr)

	for name, file := range files {
		profile := profilesByName.Get(name)
		
		decls := processor.New(file.Fst, file.Ast).Process()

		sortedDeclarations := declarations.New(declarations.Sort(decls))

		coverageCalculator := coverage.New(profile, sortedDeclarations)
		results := coverageCalculator.Process()


	}

	return nil
}
