package v0

import (
	"io"
	"os"
)

const (
	MarkdownFileName = "coverage_report.md"
	MarkdownTemplate = `
# GoMin Coverage Report {{ if .Valid }}🟢{{ else }}🔴{{ end }}

## Global Coverage
| Statements Covered  | Total Statements | Ratio |
|---|---|---|
|{{ range .Global }}{{.}}|{{ end }}

## Coverage Errors

| Place | Error |
|---|---|
{{ range .Errors }}|{{ range . }}{{.}}|{{ end }}
{{end}}

## Package Coverage
| Package Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|
{{ range .Package }}|{{ range . }}{{.}}|{{ end }}
{{end}}

## File Coverage

| Package Name | File Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|---|
{{ range .File }}|{{ range . }}{{ if . }}{{.}}{{else}}n/a{{end}}|{{ end }}
{{end}}

## Block Coverage

| Package Name | File Name | Function Name | Statements Covered  | Total Statements | Ratio | Valid | Error |
|---|---|---|---|---|---|---|---|
{{ range .Block }}|{{ range . }}{{.}}|{{ end }}
{{end}}
`
)

type Global struct {
	Covered    string
	Statements string
	Ratio      string
}

type MarkDownTemplate struct {
	Global  []string
	Errors  [][]string
	Package [][]string
	File    [][]string
	Block   [][]string
	Valid   bool
}

type GenericWriter[T any] interface {
	Write(t T) error
}

type MarkdownOutputter struct {
	WriterFactory  WriterFactory
	Format         func(node StatementNode, depth Depth) [][]string
	Validate       func(node StatementNode, place ...string) []Error
	ErrorsToRecord func(errs []Error) [][]string
}

func NewMarkdownOutputter(options ...func(*MarkdownOutputter)) *MarkdownOutputter {
	outputter := &MarkdownOutputter{
		WriterFactory:  NewWriterFactory[MarkDownTemplate](),
		Format:         Format,
		Validate:       Validate,
		ErrorsToRecord: ErrorsToRecord,
	}

	for _, option := range options {
		option(outputter)
	}

	return outputter
}

func (outputter *MarkdownOutputter) WriteOut(node StatementNode) error {
	writer, err := outputter.WriterFactory.NewWriter()
	if err != nil {
		return err
	}

	errs := outputter.ErrorsToRecord(outputter.Validate(node))
	return writer.Write(MarkDownTemplate{
		Global:  createRecord(node),
		Errors:  errs,
		Package: outputter.Format(node, PackageDepth),
		File:    outputter.Format(node, FileDepth),
		Block:   outputter.Format(node, BlockDepth),
		Valid:   len(errs) == 0,
	})
}

type FileFactory interface {
	NewFile() (io.WriteCloser, error)
}

type TemplateFactory interface {
	NewExecutor() (TemplateExecuter, error)
}

type TemplateExecutorFactory struct {
	Template       string
	CreateExecutor func(content string) (*Executor, error)
}

func NewTemplateExecutorFactory(options ...func(*TemplateExecutorFactory)) *TemplateExecutorFactory {
	factory := &TemplateExecutorFactory{
		Template:       MarkdownTemplate,
		CreateExecutor: NewExecutor,
	}

	for _, option := range options {
		option(factory)
	}

	return factory
}

func (factory *TemplateExecutorFactory) NewExecutor() (TemplateExecuter, error) {
	return factory.CreateExecutor(factory.Template)
}

type WriterFactory interface {
	NewWriter() (GenericWriter[MarkDownTemplate], error)
}

type TemplateWriterFactory[T any] struct {
	TemplateFactory TemplateFactory
	FileFactory     FileFactory
	CreateWriter    func(Exec TemplateExecuter, File io.WriteCloser) *TemplateWriter[T]
}

func NewWriterFactory[T any](options ...func(*TemplateWriterFactory[T])) *TemplateWriterFactory[T] {
	factory := &TemplateWriterFactory[T]{
		TemplateFactory: NewTemplateExecutorFactory(),
		FileFactory:     NewOSFileFactory(),
		CreateWriter:    NewTemplateWriter[T],
	}

	for _, option := range options {
		option(factory)
	}

	return factory
}

func (factory *TemplateWriterFactory[T]) NewWriter() (GenericWriter[T], error) {
	template, err := factory.TemplateFactory.NewExecutor()
	if err != nil {
		return nil, err
	}

	file, err := factory.FileFactory.NewFile()
	if err != nil {
		return nil, err
	}

	return factory.CreateWriter(template, file), nil
}

type OSFileFactory struct {
	FileName string
	Open     func(fileName string) (*os.File, error)
}

func NewOSFileFactory(options ...func(*OSFileFactory)) *OSFileFactory {
	factory := &OSFileFactory{
		FileName: MarkdownFileName,
		Open:     os.Create,
	}

	for _, option := range options {
		option(factory)
	}

	return factory
}

func (factory *OSFileFactory) NewFile() (io.WriteCloser, error) {
	return factory.Open(factory.FileName)
}

type TemplateWriter[T any] struct {
	Exec TemplateExecuter
	File io.WriteCloser
}

func NewTemplateWriter[T any](
	Exec TemplateExecuter,
	File io.WriteCloser,
) *TemplateWriter[T] {
	return &TemplateWriter[T]{
		Exec: Exec,
		File: File,
	}
}

func (writer *TemplateWriter[T]) Write(t T) error {
	return writer.Exec.Execute(writer.File, t)
}
