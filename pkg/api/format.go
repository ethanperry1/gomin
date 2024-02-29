package api

import "fmt"

type Depth int

const (
	PackageDepth Depth = iota + 1
	FileDepth
	BlockDepth
)

func Format(node StatementNode, depth Depth) [][]string {
	return format(node, 0, depth)
}

func format(node StatementNode, currentDepth, maxDepth Depth, names ...any) [][]string {
	if currentDepth == maxDepth {
		return [][]string{createRecord(node, names...)}
	}

	children := node.Children()
	if children == nil {
		return [][]string{}
	}

	var records [][]string
	for name, child := range children {
		records = append(records, format(child, currentDepth+1, maxDepth, append(names, name)...)...)
	}

	return records
}

func createRecord(node StatementNode, names ...any) []string {
	printedNames := make([]string, len(names))
	for idx, name := range names {
		printedNames[idx] = nameToString(name)
	}
	
	var reason string
	valid := node.Valid() == nil
	if !valid {
		reason = node.Valid().Error()
	}

	return append(
		printedNames,
		fmt.Sprintf("%d", node.Covered()),
		fmt.Sprintf("%d", node.Total()),
		fmt.Sprintf("%0.2f", Ratio(node)),
		fmt.Sprintf("%t", valid),
		reason,
	)
}

func nameToString(name any) string {
	switch n := name.(type) {
	case string:
		return n
	case int:
		return fmt.Sprintf("literal (index %d)", n)
	case NamePair:
		return fmt.Sprintf("%s.%s", n.X, n.Y)
	default:
		return ""
	}
}
