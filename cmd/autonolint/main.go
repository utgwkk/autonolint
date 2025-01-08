package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/utgwkk/autonolint"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	var inputJson []autonolint.Issue
	if err := json.Unmarshal(input, &inputJson); err != nil {
		log.Fatal(err)
	}

	if err := autonolint.Process(inputJson); err != nil {
		log.Fatal(err)
	}
}
