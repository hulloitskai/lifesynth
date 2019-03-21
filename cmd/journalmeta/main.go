package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stevenxie/lifesynth/journal/markdown"
	schemav1 "github.com/stevenxie/lifesynth/journal/schema/v1"
	ess "github.com/unixpickle/essentials"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		showHelp()
	}
	switch args[0] {
	case "--help", "-h":
		showHelp()
	}

	// Read file.
	data, err := ioutil.ReadFile(args[0])
	if err != nil {
		ess.Die("Reading file:", err)
	}

	// Process out metadata.
	doc, err := markdown.Process(data)
	if err != nil {
		ess.Die("Processing Markdown:", err)
	}
	data = doc.Body

	// Parse file.
	parser := schemav1.NewParser()
	ms, err := parser.ParseMeta(data)
	if err != nil {
		ess.Die("Parsing metadata from document:", err)
	}

	// Output metadata.
	for _, meta := range ms {
		fmt.Printf("• Time:     %s\n", meta.Time)
		fmt.Printf("  Location: %s\n", meta.Location)
		fmt.Printf("  Valence:  %.2f\n", meta.Valence)
		if len(meta.Tags) > 0 {
			fmt.Printf("  Tags:     %s\n", meta.Tags)
		}
	}
}

func showHelp() {
	fmt.Fprintf(os.Stderr, "Usage: %s <markdown file>\n", os.Args[0])
	os.Exit(1)
}
