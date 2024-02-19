package parser

import "fmt"

type UnknownCommandError struct {
	command string
}

func (err *UnknownCommandError) Error() string {
	return fmt.Sprintf("the command %q is unknown", err.command)
}