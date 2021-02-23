package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Base struct {
	tagData
}

func NewBase(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("base", true, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &Base{td}, nil
}

func (t *Base) Validate() error {
	if len(t.children) != 0 {
		errCtx := t.children[0].Context()
		return errCtx.NewError("HTML Error: unexpected children")
	}

	return nil
}
