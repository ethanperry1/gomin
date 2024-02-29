package api

type RuleSet interface {
	Children(name any) RuleSet
	StatementEvaluator
}

type ruleSet struct {
	evaluators []StatementEvaluator
	rules      []*ruleSet
	matchers   []Matcher
	ruleType   ArgType
}

func AddRuleSet(r ...*ruleSet) func(*ruleSet) {
	return func(rs *ruleSet) {
		rs.rules = append(rs.rules, r...)
	}
}

func AddMatcher(m ...Matcher) func(*ruleSet) {
	return func(rs *ruleSet) {
		rs.matchers = append(rs.matchers, m...)
	}
}

func AddEvaluator(s ...StatementEvaluator) func(*ruleSet) {
	return func(rs *ruleSet) {
		rs.evaluators = append(rs.evaluators, s...)
	}
}

func NewRuleSet(
	options ...func(*ruleSet),
) *ruleSet {
	ruleSet := &ruleSet{
		rules:      []*ruleSet{},
		matchers:   []Matcher{},
		evaluators: []StatementEvaluator{},
	}

	for _, option := range options {
		option(ruleSet)
	}

	return ruleSet
}

func (set *ruleSet) match(matches any) bool {
	for _, matcher := range set.matchers {
		if matcher.Match(matches) {
			return true
		}
	}

	return false
}

func (set *ruleSet) Children(name any) RuleSet {
	var rules []*ruleSet
	var evals []StatementEvaluator
	for _, rule := range set.rules {
		if rule.match(name) {
			rules = append(rules, rule.rules...)
			evals = append(evals, rule.evaluators...)
		}
	}

	return &ruleSet{
		rules:      rules,
		evaluators: evals,
	}
}

func (set *ruleSet) Apply(statements Statements) Statements {
	for _, eval := range set.evaluators {
		statements = eval.Apply(statements)
	}

	return statements
}

type StatementEvaluator interface {
	Apply(Statements) Statements
}
