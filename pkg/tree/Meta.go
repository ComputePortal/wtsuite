package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Meta struct {
	tagData
}

func NewMeta(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("meta", true, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &Meta{td}, nil
}

func (t *Meta) Validate() error {
	if len(t.children) != 0 {
		errCtx := t.children[0].Context()
		return errCtx.NewError("HTML Error: unexpected children")
	}

	return nil
}
