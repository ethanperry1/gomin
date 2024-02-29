package api

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuleSet(t *testing.T) {
	rs := NewRuleSet(
		AddRuleSet(
			NewRuleSet(
				AddRuleSet(
					NewRuleSet(
						AddMatcher(
							NewRegexpMatcher(
								regexp.MustCompile("match").MatchString,
							),
						),
						AddEvaluator(
							NewExcludeCommand(),
						),
					),
				),
				AddMatcher(
					NewRegexpMatcher(
						regexp.MustCompile("example").MatchString,
					),
				),
				AddEvaluator(
					&MinimumCommand{
						minimum: 0.1,
					},
				),
			),
			NewRuleSet(
				AddRuleSet(
					NewRuleSet(
						AddMatcher(
							NewRegexpMatcher(
								regexp.MustCompile("world").MatchString,
							),
						),
						AddEvaluator(
							&MinimumCommand{
								minimum: 1,
							},
						),
					),
				),
				AddMatcher(
					NewRegexpMatcher(
						regexp.MustCompile("example").MatchString,
					),
				),
				AddEvaluator(
					NewExcludeCommand(),
				),
			),
			NewRuleSet(
				AddMatcher(
					NewRegexpMatcher(
						regexp.MustCompile("example2").MatchString,
					),
				),
				AddEvaluator(
					&MinimumCommand{
						minimum: 0.2,
					},
				),
			),
		),
		AddEvaluator(
			&MinimumCommand{
				minimum: 0.1,
			},
		),
		AddMatcher(NewNoopMatcher()),
	)

	ns := NewNode(
		AddNode(
			"example",
			NewNode(
				AddNode(
					"match",
					NewNode(
						AddStatement(
							NewStatements(1, 20),
						),
					),
				),
			),
		),
		AddNode(
			"nomatch",
			NewNode(
				AddStatement(
					NewStatements(0, 1),
				),
			),
		),
	)

	s := eval(rs, ns)
	require.Equal(t, 0, s.Covered())
	require.Equal(t, 1, s.Total())
}
