package main

import (
	"fmt"
	"io/ioutil"

	"github.com/stevenxie/lifesynth/journal/markdown"
	ess "github.com/unixpickle/essentials"
)

const filename = "./document.md"

func main() {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		ess.Die("Reading file:", err)
	}

	doc, err := markdown.Process(data)
	if err != nil {
		ess.Die("Parsing Markdown:", err)
	}

	fmt.Printf("--- [metadata]\n%s\n", doc.Metadata)
	fmt.Printf("--- [body]\n%s\n", doc.Body)
}
