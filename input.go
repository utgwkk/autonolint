package autonolint

import "encoding/json"

type Issue struct {
	FromLinter string

	Text string

	Pos Pos
}

type Pos struct {
	Filename string

	Offset int
	Line   int
	Column int
}

type ExecResult struct {
	Issues []Issue
}

func ParseInput(in []byte) ([]Issue, error) {
	// 1. Try to parse as []Issue
	var issues []Issue
	if err := json.Unmarshal(in, &issues); err == nil {
		return issues, nil
	}

	// 2. Try to parse as ExecResult
	var execResult ExecResult
	err := json.Unmarshal(in, &execResult)
	if err == nil {
		return execResult.Issues, nil
	}

	return nil, err
}
