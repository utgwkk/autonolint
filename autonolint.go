package autonolint

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type Issue struct {
	FromLinter string

	Pos Pos
}

type Pos struct {
	Filename string

	Offset int
	Line   int
	Column int
}

func Process(issues []Issue) error {
	for _, issue := range issues {
		if err := processIssue(issue); err != nil {
			return err
		}
	}
	return nil
}

func processIssue(issue Issue) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, issue.Pos.Filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parser.ParseFile: %w", err)
	}
	pos := token.Pos(issue.Pos.Offset + 1)
	path, _ := astutil.PathEnclosingInterval(f, pos, pos)
	if len(path) == 0 {
		return nil
	}
	p := path[0]
	ast.Print(fset, p)

	// TODO: insert //nolint:xxx comment before p

	return nil
}
