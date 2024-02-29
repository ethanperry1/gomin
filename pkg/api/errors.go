package api

import "fmt"

type InvalidMinimumArgumentError struct {
	argument float64
}

func (err *InvalidMinimumArgumentError) Error() string {
	return fmt.Sprintf("received minimum %0.2f where a value between 0 and 1 was expected", err.argument)
}

type CoverageBelowThresholdError struct {
	expected float64
	actual float64
}

func (err *CoverageBelowThresholdError) Error() string {
	return fmt.Sprintf("actual coverage was %0.2f where coverage of at least %0.2f was expected", err.actual, err.expected)
}

type InvalidCommandArgumentTypeError struct {
	argType ArgType
}

func (err *InvalidCommandArgumentTypeError) Error() string {
	return fmt.Sprintf("the argType %d is not valid", err.argType)
}

type InvalidCommandArgumentValueTypeError struct {
	expected string
}

func (err *InvalidCommandArgumentValueTypeError) Error() string {
	return fmt.Sprintf("a command argument contained an unexpected type, where type %q was expected", err.expected)
}