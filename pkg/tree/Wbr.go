package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Wbr struct {
	VisibleTagData
}

func NewWbr(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("wbr", true, attr, ctx)
	return &Wbr{visTag}, err
}

func (t *Wbr) Validate() error {
	if t.NumChildren() != 0 {
		panic("should've been caught during construction")
	}

	return nil
}
