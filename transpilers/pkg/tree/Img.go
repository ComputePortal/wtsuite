package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Img struct {
	VisibleTagData
}

func NewImg(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("img", true, attr, ctx)
	return &Img{visTag}, err
}

func (t *Img) Validate() error {
	if t.NumChildren() != 0 {
		errCtx := t.Context()
		return errCtx.NewError("HTML Error: unexpected content")
	}

	return nil
}

func (t *Img) Write(indent string, nl, tab string) string {
	return t.VisibleTagData.WriteWrappedAutoHref(indent, nl, tab)
}
