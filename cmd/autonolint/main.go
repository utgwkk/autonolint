package main

import (
	"encoding/json"
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

	var inputJson []autonolint.Issue
	if err := json.Unmarshal(input, &inputJson); err != nil {
		log.Fatal(err)
	}

	if err := autonolint.Process(inputJson, *argsCommend); err != nil {
		log.Fatal(err)
	}
}
