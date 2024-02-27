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
	fmt.Fprintln(w, "Code Path \tFunction Name \tCoverage Before Directives \tCoverage After Directives \tPercentage")
	printBlocks(w, 20, "", "", c)
	w.Flush()

	fmt.Println()
	fmt.Fprintln(w, "Package Name \tFile Name \tCoverage Before Directives \tCoverage After Directives \tPercentage")
	printFiles(w, 20, 0, "", "", c)
	w.Flush()

	fmt.Println()
	fmt.Fprintln(w, "Package Name \tCoverage Before Directives \tCoverage After Directives \tPercentage")
	printPackages(w, 20, c)
	w.Flush()
}

func printBlocks(w *tabwriter.Writer, length int, parent string, name string, c Coverage) {
	children := c.Children()
	if children == nil {
		ratio := float64(c.After().Covered()) / float64(c.After().Statements())
		if c.After().Statements() == 0 {
			ratio = 0
		}
		lines := int(ratio * float64(length))
		bar := fmt.Sprintf("%s%s", strings.Repeat("*", lines), strings.Repeat(" ", length-lines))
		stmts := fmt.Sprintf("(%d/%d)\t(%d/%d)", c.Before().Covered(), c.Before().Statements(), c.After().Covered(), c.After().Statements())
		fmt.Fprintf(w, "%s (%d,%d) \t%s \t%s \t[%s]\n", parent, c.Line(), c.Col(), name, stmts, bar)
		return
	}

	for key, child := range children {
		printBlocks(w, length, filepath.Join(parent, name), key, child)
	}
}

func printFiles(w *tabwriter.Writer, length int, depth int, pkg, file string, c Coverage) {
	if depth == 2 {
		ratio := float64(c.After().Covered()) / float64(c.After().Statements())
		if c.After().Statements() == 0 {
			ratio = 0
		}
		lines := int(ratio * float64(length))
		bar := fmt.Sprintf("%s%s", strings.Repeat("*", lines), strings.Repeat(" ", length-lines))
		stmts := fmt.Sprintf("(%d/%d)\t(%d/%d)", c.Before().Covered(), c.Before().Statements(), c.After().Covered(), c.After().Statements())
		fmt.Fprintf(w, "%s \t%s \t%s \t[%s]\n", pkg, file, stmts, bar)
		return
	}

	for key, child := range c.Children() {
		printFiles(w, length, depth+1, filepath.Join(pkg, file), key, child)
	}
}

func printPackages(w *tabwriter.Writer, length int, c Coverage) {
	for key, child := range c.Children() {
		ratio := float64(child.After().Covered()) / float64(child.After().Statements())
		if child.After().Statements() == 0 {
			ratio = 0
		}
		lines := int(ratio * float64(length))
		bar := fmt.Sprintf("%s%s", strings.Repeat("*", lines), strings.Repeat(" ", length-lines))
		stmts := fmt.Sprintf("(%d/%d)\t(%d/%d)", child.Before().Covered(), child.Before().Statements(), child.After().Covered(), child.After().Statements())
		fmt.Fprintf(w, "%s \t%s \t[%s]\n", key, stmts, bar)
	}
}