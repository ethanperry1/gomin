package v0

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type FlusherBuilder interface {
	NewWriter(output io.Writer) Flusher
}

type Writer struct {
	Builder FlusherBuilder
	Delim   string
}

func NewWriter(options ...func(*Writer)) *Writer {
	writer := &Writer{
		Builder: NewTabWriterBuilder(),
		Delim:   "\t",
	}

	for _, option := range options {
		option(writer)
	}

	return writer
}

func (writer *Writer) Write(w io.Writer, records [][]string) error {
	tabWriter := writer.Builder.NewWriter(w)
	for _, row := range records {
		fmt.Fprintln(tabWriter, strings.Join(row, writer.Delim))
	}
	return tabWriter.Flush()
}

type Flusher interface {
	Flush() error
	io.Writer
}

type TabWriterBuilder struct {
	Writer   func(output io.Writer, minwidth int, tabwidth int, padding int, padchar byte, flags uint) *tabwriter.Writer
	MinWidth int
	TabWidth int
	Padding  int
	Padchar  byte
	Flags    uint
}

func NewTabWriterBuilder(options ...func(*TabWriterBuilder)) *TabWriterBuilder {
	writer := &TabWriterBuilder{
		Writer:   tabwriter.NewWriter,
		MinWidth: 1,
		TabWidth: 1,
		Padding:  1,
		Padchar:  ' ',
		Flags:    0,
	}

	for _, option := range options {
		option(writer)
	}

	return writer
}

func (builder *TabWriterBuilder) NewWriter(output io.Writer) Flusher {
	return builder.Writer(
		output,
		builder.MinWidth,
		builder.TabWidth,
		builder.Padding,
		builder.Padchar,
		builder.Flags,
	)
}
