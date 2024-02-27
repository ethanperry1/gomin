package unit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
)

func PrettyPrint(c Coverage) {
	fmt.Printf("detailed coverage results:\n")
	prettyPrint(0, "overall", c)
}

func prettyPrint(depth int, name string, c Coverage) {
	stmts := fmt.Sprintf("(%d/%d) (%d/%d)", c.Before().Covered(), c.Before().Statements(), c.After().Covered(), c.After().Statements())
	indentedName := fmt.Sprintf("%s├──%s", strings.Repeat("|\t", depth), name)

	children := c.Children()
	if children == nil {
		fmt.Printf("%s (%d,%d): %s\n", indentedName, c.Line(), c.Col(), stmts)
		return
	}

	fmt.Printf("%s: %s\n", indentedName, stmts)

	for name, child := range children {
		prettyPrint(depth+1, name, child)
	}
}

func BarChart(c Coverage) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Code Location \tCoverage Before Directives \tCoverage After Directives \tPercentage")
	barChart(w, 20, "", "", c)
	w.Flush()
}

func barChart(w *tabwriter.Writer, length int, parent string, name string, c Coverage) {
	children := c.Children()
	if children == nil {
		ratio := float64(c.After().Covered())/float64(c.After().Statements())
		if c.After().Statements() == 0 {
			ratio = 0
		}
		lines := int(ratio * float64(length))
		bar := fmt.Sprintf("%s%s", strings.Repeat("*", lines), strings.Repeat(" ", length-lines))
		stmts := fmt.Sprintf("(%d/%d)\t(%d/%d)", c.Before().Covered(), c.Before().Statements(), c.After().Covered(), c.After().Statements())
		fmt.Fprintf(w, "%s (%d,%d):\t%s\t[%s]\n", parent, c.Line(), c.Col(), stmts, bar)
		return
	}

	for key, child := range children {
		barChart(w, length, filepath.Join(parent, name), key, child)
	}
}