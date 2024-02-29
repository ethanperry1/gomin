package v0

func CreateEvaluator(
	root string,
	profile string,
	name string,
) (*Evaluator, error) {
	profParser := NewProfileReader(
		root,
		profile,
		name,
	)

	node, err := profParser.CreateNodeTree()
	if err != nil {
		return nil, err
	}

	return NewEvaluator(node), nil
}

func NewEvaluator(node Node) *Evaluator {
	return &Evaluator{
		node: node,
	}
}

type Evaluator struct {
	node Node
}

func (evaluator *Evaluator) Evaluate(
	options ...Option,
) (StatementNode, error) {
	rs, err := ParseOptions(options...)
	if err != nil {
		return nil, err
	}

	return eval(rs, evaluator.node), nil
}

func ParseOptions(
	options ...Option,
) (*ruleSet, error) {
	var ruleSets []*ruleSet
	for _, option := range options {
		sets, err := option()
		if err != nil {
			return nil, err
		}
		ruleSets = append(ruleSets, sets...)
	}

	return NewRuleSet(
		AddRuleSet(ruleSets...),
		AddMatcher(NewNoopMatcher()),
	), nil
}

func eval(r RuleSet, n Node) StatementNode {
	leaf, ok := n.Leaf()
	if ok {
		return &statementNode{
			Statements: r.Apply(leaf),
		}
	}

	m := make(map[any]StatementNode)
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
		children:   m,
	}
}
