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
	fmt.Printf("%s├──%s: (%d/%d) (%d/%d)\n", strings.Repeat("|\t", depth), name, c.Before().Covered(), c.Before().Statements(), c.After().Covered(), c.After().Statements())

	children := c.Children()
	if children == nil {
		return
	}

	for name, child := range children {
		prettyPrint(depth+1, name, child)
	}
}