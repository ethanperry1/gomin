package api

type Evaluator struct {
}

func (evaluator *Evaluator) Evaluate() {

}

func eval(r RuleSet, n Node) StatementNode {
	leaf, ok := n.Leaf()
	if ok {
		return &statementNode{
			Statements: r.Apply(leaf),
		}
	}

	m := make(map[string]StatementNode)
	s := &statements{}
	children := n.Children()
	for name, child := range children {
		evalChild := eval(r.Children(name), child)
		s.covered += evalChild.Covered()
		s.total += evalChild.Total()
		m[name] = evalChild
	}

	return &statementNode{
		Statements: r.Apply(s),
		children: m,
	}
}
