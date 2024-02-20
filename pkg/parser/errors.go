package parser

import "fmt"

type UnknownCommandError struct {
	command   string
	directive string
}

func (err *UnknownCommandError) Error() string {
	return fmt.Sprintf("the command %q is unknown (part of directive %q)", err.command, err.directive)
}
