package api

type Matcher interface {
	MatchString(s string) bool
}

type RegexpMatcher struct {
	matchString func (s string) bool
}

func (matcher *RegexpMatcher) MatchString(s string) bool {
	return matcher.matchString(s)
}

func NewRegexpMatcher(matchString func (s string) bool) *RegexpMatcher {
	return &RegexpMatcher{
		matchString: matchString,
	}
}

type NoopMatcher struct {}

func (matcher *NoopMatcher) MatchString(s string) bool {
	return true
}

func NewNoopMatcher() *NoopMatcher {
	return &NoopMatcher{}
}