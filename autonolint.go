package autonolint

import (
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"slices"

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

type InsertComment struct {
	Filename string

	Offset int

	FromLinter string

	Reason string
}

var errSkip = errors.New("skip")

func (c *InsertComment) Comment() string {
	if c.Reason != "" {
		return "//nolint:" + c.FromLinter + " // " + c.Reason
	}
	return "//nolint:" + c.FromLinter
}

func Process(issues []Issue) error {
	rewriteByFile := map[string][]*InsertComment{}

	for _, issue := range issues {
		c, err := processIssue(issue)
		if err != nil {
			if errors.Is(err, errSkip) {
				continue
			}
			return err
		}
		rewriteByFile[c.Filename] = append(rewriteByFile[c.Filename], c)
	}

	for filename, cs := range rewriteByFile {
		slices.SortFunc(cs, func(a, b *InsertComment) int {
			return a.Offset - b.Offset
		})
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parser.ParseFile: %w", err)
		}

		for _, c := range cs {
			pos := token.Pos(c.Offset)
			x := fset.File(f.Pos())
			ls := x.LineStart(x.Line(pos))
			f.Comments = append(f.Comments, &ast.CommentGroup{
				List: []*ast.Comment{
					{
						Slash: ls,
						Text:  c.Comment(),
					},
				},
			})
		}

		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		defer out.Close()
		if err := format.Node(out, fset, f); err != nil {
			return fmt.Errorf("decorator.Fprint: %w", err)
		}
	}

	return nil
}

func processIssue(issue Issue) (*InsertComment, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, issue.Pos.Filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parser.ParseFile: %w", err)
	}
	pos := token.Pos(issue.Pos.Offset + 1)
	path, _ := astutil.PathEnclosingInterval(f, pos, pos+1)
	if len(path) == 0 {
		return nil, errSkip
	}
	var p ast.Node
	for _, p1 := range path {
		p11, isStmt := p1.(ast.Stmt)
		if isStmt {
			p = p11
			break
		}
	}
	if p == nil {
		return nil, errSkip
	}

	return &InsertComment{
		Filename:   issue.Pos.Filename,
		Offset:     issue.Pos.Offset,
		FromLinter: issue.FromLinter,
		Reason:     "test",
	}, nil
}
