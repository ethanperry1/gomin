package unit

import (
	"github.com/ethanperry1/gomin/pkg/tokens"
)

type Coverage interface {
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
	before := &Result{}
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
		after:  &Result{
			Stmts: cov.Statements(),
			Covd:  cov.Covered(),
		},
	}

	return result, nil
}

type Child struct {
	name       string
	directives []tokens.LeveledComparer
	coverage   tokens.Coverage
}

func NewChild(name string,
	coverage tokens.Coverage) *Child {
	return &Child{
		name:     name,
		coverage: coverage,
	}
}

func (block *Child) Name() string {
	return block.name
}

func (block *Child) WithDirectives(directive ...tokens.LeveledComparer) {
	block.directives = append(block.directives, directive...)
}

func (block *Child) Evaluate() (Coverage, error) {
	before := &Result{
		Stmts: block.coverage.Statements(),
		Covd:  block.coverage.Covered(),
	}

	var err error
	cov := block.coverage
	for _, directive := range block.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &BlockDirectiveError{
				err:       err,
				name:      block.name,
				directive: directive.Directive(),
			}
		}
	}

	return &FullResult{
		after: &Result{
			Stmts: cov.Statements(),
			Covd:  cov.Covered(),
		}, before: before,
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
	before *Result
	after  *Result
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

type Result struct {
	Stmts int
	Covd  int
}

func (result *Result) Statements() int {
	return result.Stmts
}

func (result *Result) Covered() int {
	return result.Covd
}