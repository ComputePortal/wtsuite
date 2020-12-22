package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Br struct {
	VisibleTagData
}

func NewBr(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("br", true, attr, ctx)
	return &Br{visTag}, err
}

func (t *Br) Validate() error {
	if t.NumChildren() != 0 {
		panic("should've been caught during construction")
	}

	return nil
}
