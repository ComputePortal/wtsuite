package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Input struct {
	VisibleTagData
}

func NewInput(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("input", true, attr, ctx)
	return &Input{visTag}, err
}

func (t *Input) Validate() error {
	if t.NumChildren() != 0 {
		errCtx := t.Context()
		return errCtx.NewError("HTML Error: unexpected content")
	}

	return nil
}
