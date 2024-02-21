package tokens

import (
	"regexp"
	"strings"
)

const (
	gobar = "gobar:"
)

var directiveRegex = regexp.MustCompile(`\/\/\s*gobar:[a-z]+`)

func Tokenize(comment string) ([]string, bool) {
	if !directiveRegex.MatchString(comment) {
		return nil, false
	}

	res := strings.Split(comment, gobar)
	directive := strings.Split(res[1], " ")[0]
	components := strings.Split(directive, ":")

	return components, true
}

func Parse(pkg bool, level Level, components ...string) (LeveledComparer, error) {
	if len(components) < 1 {
		return nil, &MissingArgument{}
	}

	if !pkg && Command(components[0]) == Pkg {
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
	default:
	}

	dir, ok := CreatorsByCommand[Command(components[0])]
	if !ok {
		return nil, &UnknownCommandError{
			command:    components[0],
			components: components,
		}
	}

	cmp, err := dir(components...)
	if err != nil {
		return nil, err
	}

	return &ComparerWithLevel{
		Comparer: cmp,
		level:    level,
	}, nil
}
