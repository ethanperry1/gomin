package api

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
	Children() map[string]StatementNode
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
	total   int,
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
	children map[string]StatementNode
}

func (statements *statementNode) Children() map[string]StatementNode {
	return statements.children
}

type comparable struct {
	Statements
	level Level
}

func (comparable *comparable) Level() Level {
	return comparable.level
}

type evaluatedStatements struct {
	Statements
	previous Statements
	err error
}

func (statements *evaluatedStatements) Previous() Statements {
	return statements.previous
}

func (statements *evaluatedStatements) Valid() error {
	return statements.err
}