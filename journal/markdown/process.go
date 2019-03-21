package markdown

import (
	"bufio"
	"bytes"

	errors "golang.org/x/xerrors"
)

// A Document is a Markdown document.
type Document struct {
	Metadata []byte // YAML formatted data
	Body     []byte // pure Markdown body
}

// Process processes a Markdown document into Metadata and Body sections.
func Process(data []byte) (*Document, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	first, err := reader.ReadString('\n')
	if err != nil {
		return nil, errors.Errorf("markdown: reading first line: %w", err)
	}
	if first != "---\n" {
		return &Document{Body: data}, nil
	}

	var (
		offset = len(first)
		row    = 1
	)
	for {
		line, err := reader.ReadString('\n')
		row++
		if err != nil {
			return nil, errors.Errorf("reading line %d: %w", row, err)
		}
		if line == "---\n" {
			break
		}
		offset += len(line)
	}

	return &Document{
		Metadata: data[4:offset:offset],
		Body:     data[offset+4:],
	}, nil
}
