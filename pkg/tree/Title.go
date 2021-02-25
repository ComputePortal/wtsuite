package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
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

func (t *Title) Content() string {
  if len(t.children) > 0 {
    if txt, ok := t.children[0].(*Text); ok {
      return txt.Value()
    }
  }

  return ""
}
