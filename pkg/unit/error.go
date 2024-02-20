package unit

import "fmt"

type BlockDirectiveError struct {
	name      string
	directive []string
	err       error
}

func (err *BlockDirectiveError) Error() string {
	return fmt.Sprintf("error on block %q, directive %v: \n\t\t\t\t--> %s", err.name, err.directive, err.err.Error())
}

type FileError struct {
	name string
	err  error
}

func (err *FileError) Error() string {
	return fmt.Sprintf("error on file %q: \n\t\t\t--> %s", err.name, err.err.Error())
}

type FileDirectiveError struct {
	name      string
	directive []string
	err       error
}

func (err *FileDirectiveError) Error() string {
	return fmt.Sprintf("error on file %q, directive %v: \n\t\t\t--> %s", err.name, err.directive, err.err.Error())
}

type PackageDirectiveError struct {
	name      string
	directive []string
	err       error
}

func (err *PackageDirectiveError) Error() string {
	return fmt.Sprintf("error on package %q, directive %v: \n\t\t--> %s", err.name, err.directive, err.err.Error())
}


type PackageError struct {
	name string
	err  error
}

func (err *PackageError) Error() string {
	return fmt.Sprintf("error on package %q: \n\t\t--> %s", err.name, err.err.Error())
}

type ProjectDirectiveError struct {
	name      string
	directive []string
	err       error
}

func (err *ProjectDirectiveError) Error() string {
	return fmt.Sprintf("error on project %q, directive %v: \n\t--> %s", err.name, err.directive, err.err.Error())
}


type ProjectError struct {
	name string
	err  error
}

func (err *ProjectError) Error() string {
	return fmt.Sprintf("error on project %q: \n\t--> %s", err.name, err.err.Error())
}