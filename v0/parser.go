package api

import "regexp"

func Min(min float64, surfaces ...CommandSurface) Option {
	return func() ([]*ruleSet, error) {
		eval, err := NewMinimumCommand(min)
		if err != nil {
			return nil, err
		}

		return parseCommandSurfaces(eval, surfaces...)
	}
}

func Exclude(surfaces ...CommandSurface) Option {
	return func() ([]*ruleSet, error) {
		return parseCommandSurfaces(NewExcludeCommand(), surfaces...)
	}
}

func parseCommandSurfaces(evaluator StatementEvaluator, surfaces ...CommandSurface) ([]*ruleSet, error) {
	ruleSets := make([]*ruleSet, len(surfaces))
	for idx, surface := range surfaces {
		ruleSet, err := parseCommandSurface([]*ruleSet{}, []StatementEvaluator{evaluator}, surface)
		if err != nil {
			return nil, err
		}
		ruleSets[idx] = ruleSet
	}

	return ruleSets, nil
}

func parseCommandSurface(childSets []*ruleSet, evaluators []StatementEvaluator, surface CommandSurface) (*ruleSet, error) {
	matcher, err := parseCommandArgs(surface.Command())
	if err != nil {
		return nil, err
	}

	parent, ok := surface.Parent()
	if !ok {
		return NewRuleSet(
			AddEvaluator(evaluators...),
			AddRuleSet(childSets...),
			AddMatcher(matcher),
		), nil
	}

	return parseCommandSurface(
		[]*ruleSet{
			NewRuleSet(
				AddEvaluator(evaluators...),
				AddRuleSet(childSets...),
				AddMatcher(matcher),
			),
		},
		[]StatementEvaluator{},
		parent,
	)
}

func parseCommandArgs(args CommandArgument) (Matcher, error) {
	switch args.Type() {
	case Any:
		return NewNoopMatcher(), nil
	case Fltr:
		values, err := parseString(args.Value())
		if err != nil {
			return nil, err
		}

		reg, err := regexp.Compile(values)
		if err != nil {
			return nil, err
		}

		return NewRegexpMatcher(reg.MatchString), nil
	case Name:
		values, err := parseString(args.Value())
		if err != nil {
			return nil, err
		}

		return NewExactMatcher(values), nil
	case Pair:
		values, err := parseNamePair(args.Value())
		if err != nil {
			return nil, err
		}

		return NewExactPairMatcher(values), nil
	case Index:
		values, err := parseInt(args.Value())
		if err != nil {
			return nil, err
		}

		return NewIndexMatcher(values), nil
	}

	return nil, &InvalidCommandArgumentTypeError{
		argType: args.Type(),
	}
}

func parseString(v any) (string, error) {
	arg, ok := v.(string)
	if !ok {
		return arg, &InvalidCommandArgumentValueTypeError{
			expected: "string",
		}
	}

	return arg, nil
}

func parseNamePair(v any) (NamePair, error) {
	arg, ok := v.(NamePair)
	if !ok {
		return arg, &InvalidCommandArgumentValueTypeError{
			expected: "NamePair",
		}
	}

	return arg, nil
}

func parseInt(v any) (int, error) {
	arg, ok := v.(int)
	if !ok {
		return arg, &InvalidCommandArgumentValueTypeError{
			expected: "int",
		}
	}

	return arg, nil
}