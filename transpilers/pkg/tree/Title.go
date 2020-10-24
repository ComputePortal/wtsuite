package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Title struct {
	tagData
}

func NewTitle(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("title", false, attr, ctx)
	if err != nil {
		return nil, err
	}

	return &Title{td}, nil
}

func (t *Title) Validate() error {
	if len(t.children) != 1 {
		errCtx := t.Context()
		return errCtx.NewError("HTML Error: expected 1 text child")
	}

	if _, ok := t.children[0].(*Text); !ok {
		errCtx := t.children[0].Context()
		return errCtx.NewError("HTML Error: expected text")
	}

	return nil
}
