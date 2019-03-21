package parseutil

import (
	"bytes"
	"strings"

	errors "golang.org/x/xerrors"
)

// A TokenParser can parse text containing tokens (singular bytes separated
// from other text by whitespace).
type TokenParser struct{ tokens []byte }

// NewTokenParser creates a new TokenParser that matches
func NewTokenParser(tokens []byte) *TokenParser {
	return &TokenParser{tokens: tokens}
}

// Parse parses the raw text input into a map of tokens to their corresponding
// values.
func (tp *TokenParser) Parse(raw []byte) (map[byte]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	if len(tp.tokens) == 0 {
		return nil, errors.New("parseutil: no tokens to match text against")
	}

	var (
		result = make(map[byte]string)

		token   byte
		builder strings.Builder
	)
	for _, field := range bytes.Fields(raw) {
		if len(field) > 1 || !bytes.Contains(tp.tokens, field) {
			if token != 0 { // only save field if there's an active token
				if builder.Len() > 0 {
					builder.WriteByte(' ')
				}
				builder.Write(field)
			}
			continue
		}

		// Field is a valid token.
		if token != 0 { // there is an active token, so write it into result
			result[token] = builder.String()
			builder.Reset()
		}
		token = field[0]
	}

	// Save remaining contents of builder.
	if builder.Len() > 0 {
		result[token] = builder.String()
	}
	return result, nil
}

// ParseString is like Parse, except it accepts text in the form of a string.
func (tp *TokenParser) ParseString(text string) (map[byte]string, error) {
	return tp.Parse([]byte(text))
}

// IsToken determines whether or not c is recognized token.
func (tp *TokenParser) IsToken(c byte) bool {
	return bytes.Contains(tp.tokens, []byte{c})
}
