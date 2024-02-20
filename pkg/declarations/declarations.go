package declarations

import (
	"sort"
)

type Decl struct {
	Line      int
	Column    int
	Name      string
	Comments  []string
	FileName  string
	Directory string
}

type Declarations struct {
	sortedDecls []*Decl
}

func New(sortedDecls []*Decl) *Declarations {
	return &Declarations{
		sortedDecls: sortedDecls,
	}
}

func (declarations *Declarations) DeclByPosition(startLine, startCol int) string {
	pos := declarations.search(declarations.sortedDecls, &Decl{
		Line:   startLine,
		Column: startCol,
	})
	if pos != 0 {
		pos -= 1
	}
	return declarations.sortedDecls[pos].Name
}

func (declarations *Declarations) search(d1 []*Decl, d2 *Decl) int {
	return sort.Search(len(d1), func(i int) bool {
		if d1[i].Line == d2.Line {
			return d1[i].Column >= d2.Column
		}

		return d1[i].Line > d2.Line
	})
}

func Sort(decls []*Decl) []*Decl {
	sort.SliceStable(decls, func(i, j int) bool {
		if decls[i].Line == decls[j].Line {
			return decls[i].Column <= decls[j].Column
		}

		return decls[i].Line < decls[j].Line
	})

	return decls
}
