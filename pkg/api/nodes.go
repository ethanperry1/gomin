package api

type Node interface {
	Children() map[string]Node
	Leaf() (Statements, bool)
}

type node struct {
	nodes      map[string]Node
	statements *statements
}

func AddNode(name string, new *node) func(*node) {
	return func(n *node) {
		n.nodes[name] = new
	}
}

func AddStatement(statements *statements) func(*node) {
	return func(n *node) {
		n.statements = statements
	}
}

func NewNode(
	options ...func(*node),
) *node {
	node := &node{
		nodes: make(map[string]Node),
	}

	for _, option := range options {
		option(node)
	}

	return node
}

func (node *node) Children() map[string]Node {
	return node.nodes
}

func (node *node) Leaf() (Statements, bool) {
	return node.statements, node.statements != nil
}