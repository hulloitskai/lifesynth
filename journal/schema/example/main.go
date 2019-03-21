package main

import (
	"fmt"
	"io/ioutil"

	"github.com/stevenxie/lifesynth/journal/schema"
	ess "github.com/unixpickle/essentials"
)

const filename = "./document.md"

func main() {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		ess.Die("Reading file:", err)
	}

	doc, err := schema.Process(data)
	if err != nil {
		ess.Die("Parsing Markdown:", err)
	}

	if doc.SchemaVersion != nil {
		fmt.Printf("--- schemaVersion: %v\n", *doc.SchemaVersion)
	}
	fmt.Printf("--- [body]\n%s\n", doc.Body)
}
