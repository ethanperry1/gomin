package tokens

import "strconv"

type (
	Level   string
	Command string
	Creator func(components ...string) (Comparer, error)

	Coverage interface {
		Covered() int
		Statements() int
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
	Default Command = "command"
)

var (
	CreatorsByCommand = map[Command]Creator{
		Min:     CreateMinimumFromCommand,
		Exclude: CreateExcludeFromCommand,
	}
)

type ComparerWithLevel struct {
	Comparer
	level Level
}

func (comparer *ComparerWithLevel) Level() Level {
	return comparer.level
}

func CreateMinimumFromCommand(components ...string) (Comparer, error) {
	if len(components) < 2 {
		return nil, &MissingArgument{}
	}

	res, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return nil, &InvalidMinimumArgument{
			argument: components[1],
		}
	}

	if res < 0 || res > 1 {
		return nil, &InvalidMinimumArgument{
			argument: components[1],
		}
	}

	return &MinimumCommand{
		minimum:   res,
		directive: components,
	}, nil
}

type MinimumCommand struct {
	minimum   float64
	directive []string
}

func NewMinimum(minimum float64) *MinimumCommand {
	return &MinimumCommand{
		minimum: minimum,
	}
}


func (minimum *MinimumCommand) Compare(cov Coverage) (Coverage, error) {
	actual := float64(cov.Covered()) / float64(cov.Statements())
	if actual < minimum.minimum {
		return nil, &ActualFunctionCoverageBeneathMinimumError{
			minimum: minimum.minimum,
			actual:  actual,
		}
	}

	return cov, nil
}

func (minimum *MinimumCommand) Type() Command {
	return Min
}

func (minimum *MinimumCommand) Directive() []string {
	return minimum.directive
}

func CreateExcludeFromCommand(components ...string) (Comparer, error) {
	return &ExcludeCommand{
		directive: components,
	}, nil
}

type ExcludeCommand struct {
	directive []string
}

func (command *ExcludeCommand) Compare(cov Coverage) (Coverage, error) {
	return &ExcludeResult{}, nil
}

func (command *ExcludeCommand) Type() Command {
	return Exclude
}

func (command *ExcludeCommand) Directive() []string {
	return command.directive
}

type ExcludeResult struct{}

func (*ExcludeResult) Statements() int {
	return 0
}

func (*ExcludeResult) Covered() int {
	return 0
}
