package a

import (
	"fmt"
	"strconv"
)

func Nilerr() error {
	i, err := strconv.ParseInt("", 10, 64)
	if err == nil {
		return err
	}
	fmt.Println(i)
	return nil
}
