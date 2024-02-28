package unit

import (
	"github.com/ethanperry1/gomin/pkg/tokens"
)

type Coverage interface {
	Line() int
	Col() int
	Before() tokens.Coverage
	After() tokens.Coverage
	Children() map[string]Coverage
}

type Unit interface {
	Name() string
	Evaluate() (Coverage, error)
}

type Parent struct {
	name       string
	children   []Unit
	directives []tokens.LeveledComparer
}

func NewParent(name string) *Parent {
	return &Parent{
		name: name,
	}
}

func (parent *Parent) Name() string {
	return parent.name
}

func (parent *Parent) Children() []Unit {
	return parent.children
}

func (parent *Parent) WithChild(unit Unit) {
	parent.children = append(parent.children, unit)
}

func (parent *Parent) WithDirectives(directive ...tokens.LeveledComparer) {
	parent.directives = append(parent.directives, directive...)
}

func (parent *Parent) Evaluate() (Coverage, error) {
	before := &Result{
		Nme: parent.name,
	}
	result := NewParentResult()
	for _, child := range parent.children {
		cov, err := child.Evaluate()
		if err != nil {
			return nil, &FileError{
				err:  err,
				name: parent.name,
			}
		}

		result.ChildrenRes[child.Name()] = cov
		before.Covd += cov.After().Covered()
		before.Stmts += cov.After().Statements()
	}

	var err error
	var cov tokens.Coverage = &Result{
		Stmts: before.Stmts,
		Covd:  before.Covd,
		Nme:   before.Nme,
	}
	for _, directive := range parent.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &FileDirectiveError{
				err:       err,
				name:      parent.name,
				directive: directive.Directive(),
			}
		}
	}

	result.FullResult = FullResult{
		before: before,
		after:  cov,
	}

	return result, nil
}

type Child struct {
	name       string
	line       int
	col        int
	directives []tokens.LeveledComparer
	coverage   tokens.Coverage
}

func NewChild(
	name string,
	line int,
	col int,
	coverage tokens.Coverage,
) *Child {
	return &Child{
		name:     name,
		line:     line,
		col:      col,
		coverage: coverage,
	}
}

func (child *Child) Name() string {
	return child.name
}

func (child *Child) WithDirectives(directive ...tokens.LeveledComparer) {
	child.directives = append(child.directives, directive...)
}

func (child *Child) Evaluate() (Coverage, error) {
	before := &Result{
		Nme:   child.name,
		Stmts: child.coverage.Statements(),
		Covd:  child.coverage.Covered(),
	}

	var err error
	var after tokens.Coverage = &Result{
		Nme:   child.name,
		Stmts: child.coverage.Statements(),
		Covd:  child.coverage.Covered(),
	}
	for _, directive := range child.directives {
		after, err = directive.Compare(after)
		if err != nil {
			return nil, &BlockDirectiveError{
				err:       err,
				name:      child.name,
				directive: directive.Directive(),
			}
		}
	}

	return &FullResult{
		line:   child.line,
		col:    child.col,
		after:  after,
		before: before,
	}, err
}

type ParentResult struct {
	ChildrenRes map[string]Coverage
	FullResult
}

func NewParentResult() *ParentResult {
	return &ParentResult{
		ChildrenRes: make(map[string]Coverage),
	}
}

func (result *ParentResult) Children() map[string]Coverage {
	return result.ChildrenRes
}

type FullResult struct {
	line   int
	col    int
	before tokens.Coverage
	after  tokens.Coverage
}

func (result *FullResult) Line() int {
	return result.line
}

func (result *FullResult) Col() int {
	return result.col
}

func (result *FullResult) Before() tokens.Coverage {
	return result.before
}

func (result *FullResult) After() tokens.Coverage {
	return result.after
}

func (result *FullResult) Children() map[string]Coverage {
	return nil
}
