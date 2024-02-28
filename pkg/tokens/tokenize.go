package tokens

import (
	"regexp"
	"strings"
)

const (
	gomin = "gomin:"
)

var directiveRegex = regexp.MustCompile(`\/\/\s*gomin:[a-z]+`)

func Tokenize(comment string) ([]string, bool) {
	if !directiveRegex.MatchString(comment) {
		return nil, false
	}

	res := strings.Split(comment, gomin)
	directive := strings.Split(res[1], " ")[0]
	components := strings.Split(directive, ":")

	return components, true
}
