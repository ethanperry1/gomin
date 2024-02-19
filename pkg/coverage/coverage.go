package coverage

import (
	"golang.org/x/tools/cover"
)

type PositionMapper interface {
	DeclByPosition(startLine, startCol int) string
}

type Result struct {
	FileResult  *FileResult
	FuncResults map[string]*DeclResult
}

type FileResult struct {
	statements        int
	coveredStatements int
}

type DeclResult struct {
	statements        int
	coveredStatements int
}

type Coverage struct {
	position PositionMapper
	profile  *cover.Profile
}

func New(profile *cover.Profile, position PositionMapper) *Coverage {
	return &Coverage{
		profile:  profile,
		position: position,
	}
}

func (coverage *Coverage) Process() *Result {
	results := &Result{
		FileResult:  &FileResult{},
		FuncResults: make(map[string]*DeclResult),
	}

	for _, block := range coverage.profile.Blocks {
		decl := coverage.position.DeclByPosition(block.StartLine, block.StartCol)
		results.FuncResults[decl].statements += block.NumStmt
		results.FileResult.statements += block.NumStmt

		if block.Count > 0 {
			results.FuncResults[decl].coveredStatements += block.NumStmt
			results.FileResult.coveredStatements += block.NumStmt
		}
	}

	return results
}
