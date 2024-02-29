package api

type Matcher interface {
	Match(matches any) bool
}

type RegexpMatcher struct {
	matchString func(s string) bool
}

func (matcher *RegexpMatcher) Match(matches any) bool {
	m, ok := matches.(string)
	if !ok {
		return false
	}

	return matcher.matchString(m)
}

func NewRegexpMatcher(matchString func(s string) bool) *RegexpMatcher {
	return &RegexpMatcher{
		matchString: matchString,
	}
}

type NoopMatcher struct{}

func (matcher *NoopMatcher) Match(matches any) bool {
	return true
}

func NewNoopMatcher() *NoopMatcher {
	return &NoopMatcher{}
}

type ExactMatcher struct {
	match string
}

func NewExactMatcher(match string) *ExactMatcher {
	return &ExactMatcher{
		match: match,
	}
}

func (matcher *ExactMatcher) Match(matches any) bool {
	var res string
	switch m := matches.(type) {
	case NamePair:
		res = m.Y
	case string:
		res = m
	default:
		return false
	}

	return matcher.match == res
}

type NamePair struct {
	X string
	Y string
}

type ExactPairMatcher struct {
	pair NamePair
}

func NewExactPairMatcher(pair NamePair) *ExactPairMatcher {
	return &ExactPairMatcher{
		pair: pair,
	}
}

func (matcher *ExactPairMatcher) Match(matches any) bool {
	pair, ok := matches.(NamePair)
	if !ok {
		return false
	}
	return matcher.pair.X == pair.X && matcher.pair.Y == pair.Y
}

type IndexMatcher struct {
	index int
}

func NewIndexMatcher(index int) *IndexMatcher {
	return &IndexMatcher{
		index: index,
	}
}

func (matcher *IndexMatcher) Match(matches any) bool {
	m, ok := matches.(int)
	if !ok {
		return false
	}

	return matcher.index == m
}
