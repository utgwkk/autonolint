package indent

import "fmt"

type S struct{}

func (*S) T() {}

func F() {
	s := &S{}
	s.
		T()

	x := fmt.Sprintf("a:%s", "b")
	_ = x
}
