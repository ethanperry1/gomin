package v0

type MinimumCommand struct {
	minimum float64
}

func NewMinimumCommand(minimum float64) (*MinimumCommand, error) {
	if minimum < 0 || minimum > 1 {
		return nil, &InvalidMinimumArgumentError{}
	}

	return &MinimumCommand{
		minimum: minimum,
	}, nil
}

func (command *MinimumCommand) Apply(statements Statements) Statements {
	ratio := Ratio(statements)
	if ratio < command.minimum {
		return &evaluatedStatements{
			previous:   statements,
			Statements: statements,
			err: &CoverageBelowThresholdError{
				actual:   ratio,
				expected: command.minimum,
			},
		}
	}

	return &evaluatedStatements{
		Statements: statements,
	}
}

type ExcludeCommand struct{}

func NewExcludeCommand() *ExcludeCommand {
	return &ExcludeCommand{}
}

func (command *ExcludeCommand) Apply(old Statements) Statements {
	return &evaluatedStatements{
		previous:   old,
		Statements: &statements{},
	}
}

type NoopCommand struct{}

func NewNoopCommand() *NoopCommand {
	return &NoopCommand{}
}

func (command *NoopCommand) Apply(statements Statements) Statements {
	return &evaluatedStatements{
		Statements: statements,
	}
}
