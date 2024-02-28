package tokens

type (
	Level   string
	Command string
	Creator func(components ...string) (Comparer, error)

	Coverage interface {
		Covered() int
		Statements() int
		Name() string
	}
	LeveledComparer interface {
		Comparer
		Level() Level
	}
	Comparer interface {
		Compare(Coverage) (Coverage, error)
		Type() Command
		Directive() []string
	}
)

const (
	Package Level = "package"
	File    Level = "file"
	Block   Level = "block"
)

const (
	Min     Command = "min"
	Pkg     Command = "pkg"
	Exclude Command = "exclude"
	Regex   Command = "regex"
	Default Command = "default"
)

type ComparerWithLevel struct {
	Comparer
	level Level
}

func (comparer *ComparerWithLevel) Level() Level {
	return comparer.level
}
