package autonolint

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

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

	Line int

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
		buf := &bytes.Buffer{}
		rewritesByLine := map[int][]*InsertComment{}
		for _, c := range cs {
			rewritesByLine[c.Line] = append(rewritesByLine[c.Line], c)
		}
		in, err := os.Open(filename)
		if err != nil {
			return err
		}
		s := bufio.NewScanner(in)
		lineno := 1
		for s.Scan() {
			if cs, ok := rewritesByLine[lineno]; ok {
				for _, c := range cs {
					buf.WriteString(c.Comment() + "\n")
				}
			}
			buf.WriteString(s.Text() + "\n")
			lineno++
		}
		in.Close()

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("format.Source: %w (%s)", err, buf.String())
		}

		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("os.Create: %w", err)
		}
		defer out.Close()
		if _, err := out.Write(formatted); err != nil {
			return err
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
	line := fset.File(1).Line(p.Pos())

	return &InsertComment{
		Filename:   issue.Pos.Filename,
		Line:       line,
		FromLinter: issue.FromLinter,
		Reason:     "test",
	}, nil
}
