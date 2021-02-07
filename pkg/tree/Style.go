package tree

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Style struct {
	content string // can't be empty
	tagData
}

func NewStyle(attr *tokens.StringDict, content string, ctx context.Context) (Tag, error) {
	td, err := newTag("style", false, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &Style{content, td}, nil
}

func (t *Style) NumChildren() int {
	panic("not available")
}

func (t *Style) AppendChild(child Tag) {
	panic("not available")
}

func (t *Style) InsertChild(i int, child Tag) error {
	panic("not available")
}

func (t *Style) DeleteChild(i int) error {
	panic("not available")
}

func (t *Style) DeleteAllChildren() error {
	panic("not available")
}

func (t *Style) CollectIDs(IDMap) error {
	return nil
}

func (t *Style) Validate() error {
	if t.content == "" {
		errCtx := t.Context()
		err := errCtx.NewError("Error: expected inline content")

		panic(err)
		return err
	}

	return nil
}

func (t *Style) Write(indent string, nl, tab string) string {
	var b strings.Builder

	b.WriteString(t.writeStart(false, indent))

	b.WriteString(t.content)

	b.WriteString(t.writeStop(false, indent))

	return b.String()
}

// implement SVGTag interface
func (t *Style) Compress(vb SVGViewBox) {
}

func (t *Style) Minify() bool {
	return false
}
