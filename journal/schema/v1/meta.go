package schemav1

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"

	errors "golang.org/x/xerrors"

	"github.com/russross/blackfriday"
	"github.com/stevenxie/lifesynth/journal/schema/parseutil"
)

// Meta contains metadata that adds context to a journal entry.
type Meta struct {
	Time     string   `mapstructure:"time"`
	Location string   `mapstructure:"location"`
	Valence  float32  `mapstructure:"valence"`
	Tags     []string `mapstructure:"tags"`
}

// A set of byte tokens that corresponding to various meta fields.
const (
	metaTokenTime     byte = '>'
	metaTokenValence  byte = '!'
	metaTokenLocation byte = '@'
	metaTokenTag      byte = '#'
)

// metaTokens are byte tokens that indicate the presence of a meta field.
var metaTokens = []byte{
	metaTokenTime,
	metaTokenValence,
	metaTokenLocation,
	metaTokenTag,
}

// ParseMeta parses the meta blocks from a Markdown journal entry.
func (p *Parser) ParseMeta(entry []byte) ([]*Meta, error) {
	node := p.buildMarkdownParser().Parse(entry)
	if node == nil {
		return nil, errors.New("schemav1: parsed node was nil")
	}
	node.FirstChild = node.FirstChild.Next // skip title node

	var (
		parser = parseutil.NewTokenParser(metaTokens)
		result []*Meta

		meta Meta
		// tagset [][]string // TODO: perform levelled tag parsing.
		err error
	)
	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering {
			return blackfriday.GoToNext
		}

		metamap := make(map[byte]string)
		switch n.Type {
		case blackfriday.Code:
			if !parser.IsToken(n.Literal[0]) { // not a meta block, so skip
				return blackfriday.GoToNext
			}

			// Parse meta fields.
			if metamap, err = parser.Parse(n.Literal); err != nil {
				err = errors.Errorf("schemav1: parsing meta fields from code block: %w",
					err)
				return blackfriday.Terminate
			}

		case blackfriday.CodeBlock:
			// Skip if code block is not tagged with 'meta'.
			if strings.TrimSpace(string(n.CodeBlockData.Info)) != "meta" {
				return blackfriday.GoToNext
			}

			// Parse meta fields.
			scanner := bufio.NewScanner(bytes.NewReader(n.Literal))
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) < 2 {
					continue
				}
				if line[1] != ' ' {
					err = errors.New("schemav1: malformed line in meta block: token " +
						"must be followed by a space")
					break
				}
				metamap[line[0]] = strings.TrimSpace(line[2:])
			}
			if err == nil {
				if serr := scanner.Err(); serr != nil {
					err = errors.Errorf("schemav1: scanning meta block: %w", serr)
				}
			}
			if err != nil {
				return blackfriday.Terminate
			}

		default:
			return blackfriday.GoToNext
		}

		// Decode metamap into meta.
		if err = decodeMetaTokenMap(metamap, &meta); err != nil {
			err = errors.Errorf("schemav1: decoding meta token map: %w", err)
			return blackfriday.Terminate
		}

		// if len(meta.Tags) > 0 {
		// 	tagset = append(tagset, meta.Tags)
		// }

		// Copy meta into result.
		copy := meta
		result = append(result, &copy)
		return blackfriday.GoToNext
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ParseMetaInto is like ParseMeta, but marshals the data into a foreign value
// using gob.
func (p *Parser) ParseMetaInto(entry []byte, v interface{}) error {
	meta, err := p.ParseMeta(entry)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = gob.NewEncoder(&buf).Encode(meta); err != nil {
		return errors.Errorf("schemav1: encoding meta into gob: %w", err)
	}
	if err = gob.NewDecoder(&buf).Decode(v); err != nil {
		return errors.Errorf("schemav1: decoding gob into v: %w", err)
	}
	return nil
}

func decodeMetaTokenMap(src map[byte]string, dst *Meta) error {
	for token, val := range src {
		switch token {
		case metaTokenTime:
			dst.Time = val
		case metaTokenLocation:
			dst.Location = val

		case metaTokenTag:
			dst.Tags = strings.Fields(val)

		case metaTokenValence:
			// Check if value is a fraction.
			if components := strings.Split(val, "/"); len(components) == 2 {
				var (
					parts = make([]int, 2)
					err   error
				)
				for i, component := range components {
					if parts[i], err = strconv.Atoi(component); err != nil {
						return errors.Errorf("parsing valence fraction (segment %d): %w",
							i+1, err)
					}
				}
				dst.Valence = float32(parts[0]) / float32(parts[1])
				continue
			}

			// Value is a decimal.
			v64, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return errors.Errorf("parsing valence float: %w", err)
			}
			dst.Valence = float32(v64)
		}
	}
	return nil
}
