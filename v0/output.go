package v0

import (
	"io"
	"os"
)

const (
	BlockCovTableFile   = "block_coverage_table.txt"
	FileCovTableFile    = "file_coverage_table.txt"
	PackageCovTableFile = "package_coverage_table.txt"
	ErrorsFile          = "coverage_errors.txt"
)

type Outputter interface {
	WriteOut(node StatementNode) error
}

type RecordOutputter interface {
	WriteOut(name string, node StatementNode, depth Depth) error
}

type ErrorOutputter interface {
	WriteOut(name string, node StatementNode) error
}

type MultiFileOutputter struct {
	ErrorsFile          string
	BlockCovTableFile   string
	FileCovTableFile    string
	PackageCovTableFile string

	RecordsOut RecordOutputter
	ErrorsOut  ErrorOutputter
}

func NewMultiFileOutputter(options ...func(*MultiFileOutputter)) *MultiFileOutputter {
	outPutter := &MultiFileOutputter{
		ErrorsFile:          ErrorsFile,
		BlockCovTableFile:   BlockCovTableFile,
		FileCovTableFile:    FileCovTableFile,
		PackageCovTableFile: PackageCovTableFile,
		RecordsOut:          NewRecordOutputter(),
		ErrorsOut:           NewErrorOutputter(),
	}

	for _, option := range options {
		option(outPutter)
	}

	return outPutter
}

func (outputter *MultiFileOutputter) WriteOut(node StatementNode) error {
	err := outputter.RecordsOut.WriteOut(outputter.PackageCovTableFile, node, PackageDepth)
	if err != nil {
		return err
	}

	err = outputter.RecordsOut.WriteOut(outputter.FileCovTableFile, node, FileDepth)
	if err != nil {
		return err
	}

	err = outputter.RecordsOut.WriteOut(outputter.BlockCovTableFile, node, BlockDepth)
	if err != nil {
		return err
	}

	return outputter.ErrorsOut.WriteOut(outputter.ErrorsFile, node)
}

type RecordWriter interface {
	Write(w io.Writer, records [][]string) error
}

type RecordFileOutputter struct {
	Writer RecordWriter
	Open   func(name string) (io.WriteCloser, error)
	Format func(node StatementNode, depth Depth) [][]string
}

func NewRecordOutputter(options ...func(*RecordFileOutputter)) *RecordFileOutputter {
	outputter := &RecordFileOutputter{
		Writer: NewWriter(),
		Open:   Open,
		Format: Format,
	}

	for _, option := range options {
		option(outputter)
	}

	return outputter
}

func (output *RecordFileOutputter) WriteOut(name string, node StatementNode, depth Depth) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}

	err = output.Writer.Write(file, output.Format(node, depth))
	if err != nil {
		return err
	}

	return file.Close()
}

type ErrorFileOutputter struct {
	Writer         RecordWriter
	Open           func(name string) (io.WriteCloser, error)
	Validate       func(node StatementNode, place ...string) []Error
	ErrorsToRecord func(errs []Error) [][]string
}

func NewErrorOutputter(options ...func(*ErrorFileOutputter)) *ErrorFileOutputter {
	outputter := &ErrorFileOutputter{
		Writer:         NewWriter(),
		Open:           Open,
		Validate:       Validate,
		ErrorsToRecord: ErrorsToRecord,
	}

	for _, option := range options {
		option(outputter)
	}

	return outputter
}

func (output *ErrorFileOutputter) WriteOut(name string, node StatementNode) error {
	errs := output.Validate(node, "doe")

	file, err := output.Open(name)
	if err != nil {
		return err
	}

	err = output.Writer.Write(file, output.ErrorsToRecord(errs))
	if err != nil {
		return err
	}

	return file.Close()
}

func Open(name string) (io.WriteCloser, error) {
	return os.Create(name)
}