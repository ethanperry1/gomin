package parser

import (
	"gobar/pkg/declarations"
	"gobar/pkg/tokens"
	"regexp"
	"strings"
)

const (
	gobar = "gobar:"
)

var directiveRegex = regexp.MustCompile(`\/\/\s*gobar:[a-z]+`)

type Parser struct {
	commands map[tokens.Command]tokens.Creator
}

func New(commands map[tokens.Command]tokens.Creator) *Parser {
	return &Parser{
		commands: commands,
	}
}

func (parser *Parser) ParseComment(comment string, parents ...tokens.Comparer) (tokens.Comparer, error) {
	if !directiveRegex.MatchString(comment) {
		return nil, nil
	}

	res := strings.Split(comment, gobar)
	directive := strings.Split(res[1], " ")[0]
	components := strings.Split(directive, ":")

	dir, ok := parser.commands[tokens.Command(components[0])]
	if !ok {
		return nil, &UnknownCommandError{
			command: components[0],
		}
	}

	return dir(components...)
}

func (parser *Parser) ParseComments(comments []string, parents ...tokens.Comparer) ([]tokens.Comparer, error) {
	var comparers []tokens.Comparer
	for _, comment := range comments {
		command, err := parser.ParseComment(comment, parents...)
		if err != nil {
			return nil, err
		}

		if command != nil {
			comparers = append(comparers, command)
		}
	}

	return comparers, nil
}

func (parser *Parser) ParseFile(fileComments []string, decls []*declarations.Decl) ([]tokens.Comparer, error){
	fileComparers, err := parser.ParseComments(fileComments)
	if err != nil {
		return nil, err
	}

	for _, decl := range decls {
		fileComparers, err := parser.ParseComments(fileComments)
		if err != nil {
			return nil, err
		}
	}
}