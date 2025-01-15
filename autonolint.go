package autonolint

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"

	"github.com/dave/dst/decorator"
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
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parser.ParseFile: %w", err)
		}
		dec := decorator.NewDecorator(fset)

		for _, c := range cs {
			pos := token.Pos(c.Offset)
			path, _ := astutil.PathEnclosingInterval(f, pos, pos)
			n, err := dec.DecorateNode(path[0])
			if err != nil {
				return fmt.Errorf("dec.DecorateNode: %w", err)
			}

			n.Decorations().End.Prepend(c.Comment())
		}

		ff, err := dec.DecorateFile(f)
		if err != nil {
			return fmt.Errorf("decorator.DecorateFile: %w", err)
		}
		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		defer out.Close()
		if err := decorator.Fprint(out, ff); err != nil {
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
	path, _ := astutil.PathEnclosingInterval(f, pos, pos)
	if len(path) == 0 {
		return nil, errSkip
	}
	p := path[0]

	return &InsertComment{
		Filename:   issue.Pos.Filename,
		Offset:     int(p.Pos()),
		FromLinter: issue.FromLinter,
		Reason:     "",
	}, nil
}
