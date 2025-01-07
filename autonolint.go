package autonolint

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
	// TODO: implement it
	return nil
}
