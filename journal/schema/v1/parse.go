package schemav1

import (
	"github.com/russross/blackfriday"
	errors "golang.org/x/xerrors"
)

// A Parser can parse journal entries.
type Parser struct{}

// NewParser creates a new Parser.
func NewParser() *Parser { return new(Parser) }

func (*Parser) buildMarkdownParser() *blackfriday.Markdown {
	return blackfriday.New(blackfriday.WithExtensions(
		blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs,
	))
}

func parseTitle(n *blackfriday.Node) (string, error) {
	if n.Type != blackfriday.Heading {
		return "", errors.Errorf("expected first node to be a header, instead "+
			"got '%s'", n.Type)
	}
	if n.HeadingData.Level != 1 {
		return "", errors.Errorf("expected title (first header) to be an H1, "+
			"but instead got an H%d", n.HeadingData.Level)
	}

	title := n.FirstChild
	if title.Type != blackfriday.Text {
		return "", errors.Errorf("expected title (first header) to contain text, "+
			"instead found '%s'", title.Type)
	}
	return string(title.Literal), nil
}

// // Parse parses the entry data for
// func (p *Parser) Parse(data []byte) (err error) {
// 	node := p.md.Parse(data)
// 	node.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
// 		if !entering {
// 			return blackfriday.GoToNext
// 		}

// 		fmt.Println(n)
// 		return blackfriday.GoToNext
// 	})
// 	return nil
// }

// // Read implements io.Reader for Parser.
// func (p *Parser) Read(data []byte) (n int, err error) {
// 	if err := p.Parse(data); err != nil {
// 		return 0, err
// 	}
// 	return len(data), nil
// }
