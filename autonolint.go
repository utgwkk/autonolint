package autonolint

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
)

type InsertComment struct {
	Filename string

	Line int

	FromLinter string

	Reason string
}

func (c *InsertComment) Comment() string {
	if c.Reason != "" {
		return "//nolint:" + c.FromLinter + " // " + c.Reason
	}
	return "//nolint:" + c.FromLinter
}

func Process(issues []Issue, comment string) error {
	rewriteByFile := map[string][]*InsertComment{}

	for _, issue := range issues {
		c, err := processIssue(issue, comment)
		if err != nil {
			return err
		}
		rewriteByFile[c.Filename] = append(rewriteByFile[c.Filename], c)
	}

	for filename, cs := range rewriteByFile {
		rewritesByLine := map[int][]*InsertComment{}
		for _, c := range cs {
			rewritesByLine[c.Line] = append(rewritesByLine[c.Line], c)
		}
		buf := &bytes.Buffer{}
		if err := writeFileWithNolintComments(filename, buf, rewritesByLine); err != nil {
			return fmt.Errorf("writeFileWithNolintComments: %w", err)
		}

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

func processIssue(issue Issue, reason string) (*InsertComment, error) {
	fset := token.NewFileSet()
	if _, err := parser.ParseFile(fset, issue.Pos.Filename, nil, parser.ParseComments); err != nil {
		return nil, fmt.Errorf("parser.ParseFile: %w", err)
	}
	pos := token.Pos(issue.Pos.Offset)
	line := fset.File(pos).Line(pos)

	return &InsertComment{
		Filename:   issue.Pos.Filename,
		Line:       line,
		FromLinter: issue.FromLinter,
		Reason:     reason,
	}, nil
}

func writeFileWithNolintComments(filename string, out io.StringWriter, rewritesByLine map[int][]*InsertComment) error {
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()

	s := bufio.NewScanner(in)
	lineno := 1
	for s.Scan() {
		if cs, ok := rewritesByLine[lineno]; ok {
			for _, c := range cs {
				out.WriteString(c.Comment() + "\n")
			}
		}
		out.WriteString(s.Text() + "\n")
		lineno++
	}

	return nil
}
