package tokens

import "fmt"

type ActualFunctionCoverageBeneathMinimumError struct {
	minimum     float64
	actual      float64
	funcName    string
	fileName    string
	packageName string
}

func (err *ActualFunctionCoverageBeneathMinimumError) Error() string {
	return ""
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
