package schema

import (
	"github.com/stevenxie/lifesynth/journal/markdown"
	errors "golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

// A Document is a Markdown document that may contain a schema version.
type Document struct {
	SchemaVersion *int
	Body          []byte
}

// Process processes data into a Document.
func Process(data []byte) (*Document, error) {
	doc, err := markdown.Process(data)
	if err != nil {
		return nil, errors.Errorf("schema: processing Markdown: %w", err)
	}

	// No metadata present.
	if doc.Metadata == nil {
		return &Document{Body: doc.Body}, nil
	}

	// Metadata is present, so parse it.
	var meta struct {
		SchemaVersion int `yaml:"schemaVersion"`
	}
	if err := yaml.Unmarshal(doc.Metadata, &meta); err != nil {
		return nil, errors.Errorf("schema: unmarshalling metadata as YAML: %w", err)
	}
	return &Document{
		SchemaVersion: &meta.SchemaVersion,
		Body:          doc.Body,
	}, nil
}
