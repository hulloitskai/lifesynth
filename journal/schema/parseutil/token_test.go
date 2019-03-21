package parseutil_test

import (
	"testing"

	"github.com/stevenxie/lifesynth/journal/schema/parseutil"
)

func TestTokenParser_Parse(t *testing.T) {
	var (
		tokens = []byte{'>', '!', '#', '@'}
		parser = parseutil.NewTokenParser(tokens)
	)
	const text = "> 9:36 AM ! 6/10 @ Versus Coffee # morning work/uniiverse"

	// Parse text.
	result, err := parser.ParseString(text)
	if err != nil {
		t.Errorf("Failed to parse text: %v", err)
	}

	// Log results.
	for key, val := range result {
		t.Logf("result[%c] == %s", key, val)
	}
}
