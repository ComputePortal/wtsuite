package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Link struct {
	tagData
}

func NewStyleSheetLink(relPath string, ctx context.Context) (Tag, error) {
	attr := tokens.NewEmptyStringDict(ctx)

	hrefToken, err := tokens.NewString("href", ctx)
	if err != nil {
		return nil, err
	}

	srcToken, err := tokens.NewString(relPath, ctx)
	if err != nil {
		return nil, err
	}
	attr.Set(hrefToken, srcToken)

	relToken, err := tokens.NewString("rel", ctx)
	if err != nil {
		return nil, err
	}
	stylesheetToken, err := tokens.NewString("stylesheet", ctx)
	if err != nil {
		return nil, err
	}
	attr.Set(relToken, stylesheetToken)

	typeToken, err := tokens.NewString("type", ctx)
	if err != nil {
		return nil, err
	}
	textcssToken, err := tokens.NewString("text/css", ctx)
	if err != nil {
		return nil, err
	}
	attr.Set(typeToken, textcssToken)

	return NewLink(attr, ctx)
}

func NewLink(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	td, err := newTag("link", true, attr, ctx)
	if err != nil {
		return nil, err
	}
	return &Link{td}, nil
}

func (t *Link) Validate() error {
	if len(t.children) != 0 {
		errCtx := t.children[0].Context()
		return errCtx.NewError("HTML Error: unexpected children")
	}

	return nil
}
