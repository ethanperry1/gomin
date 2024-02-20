package tokens

import "strconv"

type (
	Command string
	Creator func(components ...string) (Comparer, error)
)

const (
	Min     Command = "min"
	Pkg     Command = "pkg"
	Exclude Command = "exclude"
)

var (
	CreatorsByCommand = map[Command]Creator{
		Min:     CreateMinimumFromCommand,
		Pkg:     CreatePackageCommandFromCommand,
		Exclude: CreateExcludeFromCommand,
	}
)

type Coverage interface {
	Covered() int
	Statements() int
}

type Comparer interface {
	Compare(Coverage) (Coverage, error)
	Type() Command
	Directive() []string
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

func CreatePackageCommandFromCommand(components ...string) (Comparer, error) {
	if len(components) < 2 {
		return nil, &MissingArgument{}
	}

	switch components[1] {
	case "min":
		return CreateMinimumFromCommand(components[1:]...)
	case "exclude":
		return CreateExcludeFromCommand(components[1:]...)
	default:
		return nil, &InvalidPackageCommandError{
			command: components[1],
		}
	}
}

func CreateExcludeFromCommand(components ...string) (Comparer, error) {
	return &ExcludeCommand{
		directive: components,
	}, nil
}

type ExcludeCommand struct {
	directive []string
}


func (minimum *ExcludeCommand) Compare(cov Coverage) (Coverage, error) {
	return &ExcludeResult{}, nil
}

func (minimum *ExcludeCommand) Type() Command {
	return Exclude
}

func (minimum *ExcludeCommand) Directive() []string {
	return minimum.directive
}

type ExcludeResult struct{}

func (*ExcludeResult) Statements() int {
	return 0
}

func (*ExcludeResult) Covered() int {
	return 0
}