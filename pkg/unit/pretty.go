package unit

import (
	"fmt"
	"strings"
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
