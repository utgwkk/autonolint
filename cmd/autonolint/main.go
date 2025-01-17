package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/utgwkk/autonolint"
)

var (
	argsCommend = flag.String("comment", "", "a comment after //nolint directive")
)

func main() {
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	issues, err := autonolint.ParseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	if err := autonolint.Process(issues, *argsCommend); err != nil {
		log.Fatal(err)
	}
}
