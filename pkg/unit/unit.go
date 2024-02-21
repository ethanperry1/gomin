package unit

import (
	"gobar/pkg/tokens"
)

type Coverage interface {
	Statements() int
	Covered() int
}

type Parent interface {
	Children() []Unit
	Unit
}

type Unit interface {
	Name() string
	Evaluate() (Coverage, error)
}

type Project struct {
	name       string
	children   []Unit
	directives []tokens.Comparer
}

func NewProject(name string) *Project {
	return &Project{
		name: name,
	}
}

func (project *Project) WithChild(unit Unit) {
	project.children = append(project.children, unit)
}

func (project *Project) WithDirectives(directive ...tokens.Comparer) {
	project.directives = append(project.directives, directive...)
}

func (project *Project) Name() string {
	return project.name
}

func (project *Project) Children() []Unit {
	return project.children
}

func (project *Project) Evaluate() (Coverage, error) {
	result := &Result{}
	for _, child := range project.children {
		cov, err := child.Evaluate()
		if err != nil {
			return nil, &ProjectError{
				err:  err,
				name: project.name,
			}
		}

		result.Covd += cov.Covered()
		result.Stmts += cov.Statements()
	}

	var err error
	var cov tokens.Coverage = result
	for _, directive := range project.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &ProjectDirectiveError{
				err:       err,
				directive: directive.Directive(),
				name:      project.name,
			}
		}
	}

	return cov, nil
}

type Package struct {
	name       string
	children   []Unit
	directives []tokens.LeveledComparer
}

func NewPackage(name string) *Package {
	return &Package{
		name: name,
	}
}

func (pack *Package) Name() string {
	return pack.name
}

func (pack *Package) Children() []Unit {
	return pack.children
}

func (pack *Package) WithChild(unit Unit) {
	pack.children = append(pack.children, unit)
}

func (pack *Package) WithDirectives(directive ...tokens.LeveledComparer) {
	pack.directives = append(pack.directives, directive...)
}

func (pack *Package) Evaluate() (Coverage, error) {
	result := &Result{}
	for _, child := range pack.children {
		cov, err := child.Evaluate()
		if err != nil {
			return nil, &PackageError{
				err:  err,
				name: pack.name,
			}
		}

		result.Covd += cov.Covered()
		result.Stmts += cov.Statements()
	}

	var err error
	var cov tokens.Coverage = result
	for _, directive := range pack.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &PackageDirectiveError{
				err:       err,
				directive: directive.Directive(),
				name:      pack.name,
			}
		}
	}

	return cov, nil
}

type File struct {
	name       string
	children   []Unit
	directives []tokens.LeveledComparer
}

func NewFile(name string) *File {
	return &File{
		name: name,
	}
}

func (file *File) Name() string {
	return file.name
}

func (file *File) Children() []Unit {
	return file.children
}

func (file *File) WithChild(unit Unit) {
	file.children = append(file.children, unit)
}

func (file *File) WithDirectives(directive ...tokens.LeveledComparer) {
	file.directives = append(file.directives, directive...)
}

func (file *File) Evaluate() (Coverage, error) {
	result := &Result{}
	for _, child := range file.children {
		cov, err := child.Evaluate()
		if err != nil {
			return nil, &FileError{
				err:  err,
				name: file.name,
			}
		}

		result.Covd += cov.Covered()
		result.Stmts += cov.Statements()
	}

	var err error
	var cov tokens.Coverage = result
	for _, directive := range file.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &FileDirectiveError{
				err:       err,
				name:      file.name,
				directive: directive.Directive(),
			}
		}
	}

	return cov, nil
}

type Block struct {
	name       string
	directives []tokens.LeveledComparer
	coverage   Coverage
}

func NewBlock(name string,
	coverage Coverage) *Block {
	return &Block{
		name:     name,
		coverage: coverage,
	}
}

func (block *Block) Name() string {
	return block.name
}

func (block *Block) WithDirectives(directive ...tokens.LeveledComparer) {
	block.directives = append(block.directives, directive...)
}

func (block *Block) Evaluate() (Coverage, error) {
	var err error
	cov := block.coverage
	for _, directive := range block.directives {
		cov, err = directive.Compare(cov)
		if err != nil {
			return nil, &BlockDirectiveError{
				err:       err,
				name:      block.name,
				directive: directive.Directive(),
			}
		}
	}

	return cov, err
}

type Result struct {
	Stmts int
	Covd  int
}

func (result *Result) Statements() int {
	return result.Stmts
}

func (result *Result) Covered() int {
	return result.Covd
}
