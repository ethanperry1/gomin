package evaluate

import "fmt"

type InvalidPackageDirectiveError struct {
	err  error
	name string
}

func (err *InvalidPackageDirectiveError) Error() string {
	return fmt.Sprintf("invalid directive specified in package %q: \n\t--> %s", err.name, err.err.Error())
}

type InvalidFileDirectiveError struct {
	err  error
	dir  string
	name string
}

func (err *InvalidFileDirectiveError) Error() string {
	return fmt.Sprintf("invalid directive specified in package %q, file %q: \n\t--> %s", err.dir, err.name, err.err.Error())
}

type InvalidBlockDirectiveError struct {
	err   error
	dir   string
	name  string
	block string
}

func (err *InvalidBlockDirectiveError) Error() string {
	return fmt.Sprintf("invalid directive specified in package %q, file %q, block %q: \n\t--> %s", err.dir, err.name, err.block, err.err.Error())
}
