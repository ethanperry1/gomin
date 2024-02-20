package coverage

import (
	"golang.org/x/tools/cover"
)

type PositionMapper interface {
	DeclByPosition(startLine, startCol int) string
}

type Result struct {
	Stmts int
	Covd  int
}

func (result *Result) Statements() int {
	return result.Stmts
}

func (result *Result) Covered() int {
	return result.Covd
}

type Coverage struct {
	position PositionMapper
}

func New(position PositionMapper) *Coverage {
	return &Coverage{
		position: position,
	}
}

func (coverage *Coverage) ProcessCoverage(profile *cover.Profile) map[string]*Result {
	results := make(map[string]*Result)

	for _, block := range profile.Blocks {
		decl := coverage.position.DeclByPosition(block.StartLine, block.StartCol)
		result, ok := results[decl]
		if !ok {
			result = &Result{}
			results[decl] = result
		}

		result.Stmts += block.NumStmt

		if block.Count > 0 {
			result.Covd += block.NumStmt
		}
	}

	return results
}

func Aggregate(resultList []*Result) *Result {
	aggregatedResult := &Result{}
	for _, result := range resultList {
		aggregatedResult.Covd += result.Covd
		aggregatedResult.Stmts += result.Stmts
	}

	return aggregatedResult
}
