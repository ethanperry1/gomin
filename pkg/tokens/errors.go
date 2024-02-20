package tokens

import "fmt"

type ActualFunctionCoverageBeneathMinimumError struct {
	minimum float64
	actual  float64
}

func (err *ActualFunctionCoverageBeneathMinimumError) Error() string {
	return fmt.Sprintf("actual coverage did not meet the minimum coverage bar (expected %0.2f and got %0.2f)", err.minimum, err.actual)
}

type MissingArgument struct{}

func (err *MissingArgument) Error() string {
	return "an argument for the command was expected"
}

type InvalidMinimumArgument struct {
	argument string
}

func (err *InvalidMinimumArgument) Error() string {
	return fmt.Sprintf("received argument %q where a well-formed float between 0 and 1 was expected", err.argument)
}

type InvalidPackageCommandError struct {
	command   string
}

func (err *InvalidPackageCommandError) Error() string {
	return fmt.Sprintf("the command %q is invalid or cannot be used as a package directive", err.command)
}
