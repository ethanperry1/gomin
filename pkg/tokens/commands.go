package tokens

type MinimumCommand struct {
	minimum   float64
	directive []string
}

func NewMinimum(minimum float64) *MinimumCommand {
	return &MinimumCommand{
		minimum: minimum,
	}
}

func (minimum *MinimumCommand) Compare(cov Coverage) (Coverage, error) {
	actual := float64(cov.Covered()) / float64(cov.Statements())
	if actual < minimum.minimum {
		return nil, &ActualFunctionCoverageBeneathMinimumError{
			minimum: minimum.minimum,
			actual:  actual,
		}
	}

	return cov, nil
}

func (minimum *MinimumCommand) Type() Command {
	return Min
}

func (minimum *MinimumCommand) Directive() []string {
	return minimum.directive
}

func CreateExcludeFromCommand(components ...string) (Comparer, error) {
	return &ExcludeCommand{
		directive: components,
	}, nil
}

type ExcludeCommand struct {
	directive []string
}

func (command *ExcludeCommand) Compare(cov Coverage) (Coverage, error) {
	return &ExcludeResult{}, nil
}

func (command *ExcludeCommand) Type() Command {
	return Exclude
}

func (command *ExcludeCommand) Directive() []string {
	return command.directive
}

type ExcludeResult struct{}

func (*ExcludeResult) Statements() int {
	return 0
}

func (*ExcludeResult) Covered() int {
	return 0
}

func (*ExcludeResult) Name() string {
	return ""
}

type RegexCommand struct {
	directive []string
	command   Comparer
	regex     *regexp.Regexp
}

func (command *RegexCommand) Compare(cov Coverage) (Coverage, error) {
	if command.regex.MatchString(cov.Name()) {
		return command.command.Compare(cov)
	}

	return cov, nil
}

func (command *RegexCommand) Type() Command {
	return Regex
}

func (command *RegexCommand) Directive() []string {
	return command.directive
}
