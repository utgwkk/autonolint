package autonolint

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
