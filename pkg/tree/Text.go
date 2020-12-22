package tree

import (
	"strings"

	"../tokens/context"
)

type Text struct {
	value string
	LeafTag
}

func NewTextValue(value string, ctx context.Context) Text {
	return Text{value, NewLeafTag(ctx)}
}

func NewText(value string, ctx context.Context) *Text {
	return &Text{value, NewLeafTag(ctx)}
}

func (t *Text) Value() string {
	return t.value
}

func (t *Text) Write(indent string, nl, tab string) string {
	return strings.Trim(t.value, "\n")
}

func (t *Text) Compress(vb SVGViewBox) {
	return
}

func (t *Text) Minify() bool {
	// XXX: can empty text be minified?
	return false
}
