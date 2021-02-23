package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Col struct {
	VisibleTagData
}

func NewCol(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("col", true, attr, ctx)
	return &Col{visTag}, err
}

func (t *Col) Validate() error {
	if t.NumChildren() != 0 {
		panic("should've been caught during construction")
	}

	return nil
}
