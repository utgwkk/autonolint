package autonolint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
)

func TestAutonolintProcess(t *testing.T) {
	if _, err := exec.LookPath("golangci-lint"); err != nil {
		t.Fatal(err)
	}

	inDir := os.DirFS("testdata")
	outDir, err := os.MkdirTemp(path.Join("testoutput"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if t.Failed() {
			t.Log("process result remains in ", outDir)
		} else {
			err := os.RemoveAll(outDir)
			if err != nil {
				t.Log("failed to remove dir:", err)
			}
		}
	})
	if err := os.CopyFS(outDir, inDir); err != nil {
		t.Fatal(err)
	}

	dirs, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, dir := range dirs {
		t.Run(dir.Name(), func(t *testing.T) {
			t.Parallel()

			targetDir := path.Join(outDir, dir.Name())
			decoded, err := executeGolangciLint(targetDir)
			if err != nil {
				t.Fatal(err)
			}

			if err := Process(decoded.Issues, "test"); err != nil {
				t.Fatal(err)
			}

			decoded2, err := executeGolangciLint(targetDir)
			if err != nil {
				t.Fatal(err)
			}
			if issues := decoded2.Issues; len(issues) > 0 {
				t.Fail()
				t.Log("Found following issues:")
				for _, issue := range issues {
					filename, line, column := issue.Pos.Filename, issue.Pos.Line, issue.Pos.Column
					t.Logf("\t%s:%d:%d: %s (%s)", filename, line, column, issue.Text, issue.FromLinter)
				}
			}
		})
	}
}

func executeGolangciLint(dir string) (*ExecResult, error) {
	lintOut := &bytes.Buffer{}

	configPath := path.Join(dir, ".golangci.yml")
	cmd := exec.Command("golangci-lint", "run", "--out-format=json", "--config="+configPath, dir)
	cmd.Stdout = lintOut
	_ = cmd.Run()

	var decoded ExecResult
	if err := json.NewDecoder(lintOut).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}
	return &decoded, nil
}
