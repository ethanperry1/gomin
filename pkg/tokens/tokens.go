package tokens

import "strconv"

type (
	Command string
	Creator func(components ...string) (Comparer, error)
)

const (
	Min     Command = "min"
	Exclude Command = "exclude"
)

var (
	CreatorsByCommand = map[Command]Creator{
		Min:     CreateMinimumFromCommand,
		Exclude: nil,
	}
)

type Coverage interface {
	Covered() int
	Statements() int
}

type Comparer interface {
	Compare(Coverage) (Coverage, error)
	Type() Command
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
		minimum: res,
	}, nil
}

type MinimumCommand struct {
	minimum     float64
	funcName    string
	fileName    string
	packageName string
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
			minimum:     minimum.minimum,
			actual:      actual,
			funcName:    minimum.funcName,
			fileName:    minimum.fileName,
			packageName: minimum.packageName,
		}
	}

	return cov, nil
}

func (minimum *MinimumCommand) Type() Command {
	return Min
}
