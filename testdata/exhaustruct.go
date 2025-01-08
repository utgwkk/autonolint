package a

import "fmt"

type S struct {
	A int
	B string
}

func Exhaustruct() {
	s := &S{
		A: 1,
	}
	fmt.Println(s)
}
