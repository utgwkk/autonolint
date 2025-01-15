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
	argsReason = flag.String("reason", "", "The reason for //nolint")
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

	if err := autonolint.Process(inputJson, *argsReason); err != nil {
		log.Fatal(err)
	}
}
