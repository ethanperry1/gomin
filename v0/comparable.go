package v0

import (
	"fmt"
	"strings"
)

type Level int

const (
	Proj Level = iota
	Pack
	File
	Bloc
)

type Comparable interface {
	Level
	Statements
}

type StatementNode interface {
	Statements
	Children() map[any]StatementNode
}

type Statements interface {
	Covered() int
	Total() int
	Valid() error
	Previous() Statements
}

type statements struct {
	covered int
	total   int
}

func NewStatements(
	covered int,
	total int,
) *statements {
	return &statements{
		covered: covered,
		total:   total,
	}
}

func (statements *statements) Previous() Statements {
	return nil
}

func (statements *statements) Valid() error {
	return nil
}

func (statements *statements) Covered() int {
	return statements.covered
}

func (statements *statements) Total() int {
	return statements.total
}

func Ratio(statements Statements) float64 {
	return float64(statements.Covered()) / float64(statements.Total())
}

type statementNode struct {
	Statements
	children map[any]StatementNode
}

func (statements *statementNode) Children() map[any]StatementNode {
	return statements.children
}

type evaluatedStatements struct {
	Statements
	previous Statements
	err      error
}

func (statements *evaluatedStatements) Previous() Statements {
	return statements.previous
}

func (statements *evaluatedStatements) Valid() error {
	return statements.err
}

type Error struct {
	Err   error
	Place []string
}

func Validate(node StatementNode, place ...string) []Error {
	errs := ValidatePreviousStatements(node, place)

	children := node.Children()
	for name, child := range children {
		errs = append(errs, Validate(child, append(place, nameToString(name))...)...)
	}

	return errs
}

func ValidatePreviousStatements(statements Statements, place []string) []Error {
	if statements == nil {
		return []Error{}
	}

	var errs []Error
	if statements.Valid() != nil {
		errs = []Error{
			{
				Err:   statements.Valid(),
				Place: place,
			},
		}
	}

	return append(errs, ValidatePreviousStatements(statements.Previous(), place)...)
}

func ErrorsToRecord(errs []Error) [][]string {
	records := make([][]string, len(errs))
	for idx, err := range errs {
		records[idx] = []string{fmt.Sprintf("[%s]", strings.Join(err.Place, ", ")), err.Err.Error()}
	}
	return records
}