package tokens

import (
	"regexp"
	"strconv"
)

func Parse(pkg bool, level Level, components ...string) (LeveledComparer, error) {
	if len(components) < 1 {
		return nil, &MissingArgument{}
	}

	if !pkg == (Command(components[0]) == Pkg) {
		return nil, nil
	}

	switch Command(components[0]) {
	case Pkg:
		return Parse(false, level, components[1:]...)
	case Default:
		switch l := Level(components[1]); l {
		case File:
			return Parse(false, File, components[2:]...)
		case Block:
			return Parse(false, Block, components[2:]...)
		default:
			return nil, &InvalidDefaultCommandError{
				argument: l,
			}
		}
	}

	cmp, err := ParseCommand(components...)
	if err != nil {
		return nil, err
	}
	return &ComparerWithLevel{
		Comparer: cmp,
		level:    level,
	}, nil
}

func ParseCommand(components ...string) (Comparer, error) {
	switch Command(components[0]) {
	case Min:
		return ParseMin(components...)
	case Exclude:
		return ParseExclude(components...)
	case Regex:
		return ParseRegex(components...)
	}

	return nil, &UnknownCommandError{
		command:    components[0],
		components: components,
	}
}

func ParseRegex(components ...string) (Comparer, error) {
	if len(components) < 3 {
		return nil, &MissingArgument{}
	}

	reg, err := regexp.Compile(components[1])
	if err != nil {
		return nil, &InvalidRegexError{
			reg: components[1],
			err: err,
		}
	}

	subCommand, err := ParseCommand(components[2:]...)
	if err != nil {
		return nil, err
	}

	return &RegexCommand{
		directive: components,
		command:   subCommand,
		regex:     reg,
	}, nil
}

func ParseExclude(components ...string) (Comparer, error) {
	return &ExcludeCommand{
		directive: components,
	}, nil
}

func ParseMin(components ...string) (Comparer, error) {
	if len(components) < 2 {
		return nil, &MissingArgument{}
	}

	res, err := strconv.ParseFloat(components[1], 64)
	if err != nil {
		return nil, &InvalidMinimumArgument{
			argument: components[1],
		}
	}

	if res < 0 || res > 1 {
		return nil, &InvalidMinimumArgument{
			argument: components[1],
		}
	}

	return &MinimumCommand{
		minimum:   res,
		directive: components,
	}, nil
}
