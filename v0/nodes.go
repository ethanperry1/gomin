package v0

type Node interface {
	Children() map[any]Node
	Leaf() (Statements, bool)
}

type node struct {
	nodes      map[any]Node
	statements *statements
}

func AddNode(name any, new *node) func(*node) {
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
		nodes: make(map[any]Node),
	}

	for _, option := range options {
		option(node)
	}

	return node
}

func (node *node) Children() map[any]Node {
	return node.nodes
}

func (node *node) Leaf() (Statements, bool) {
	return node.statements, node.statements != nil
}
