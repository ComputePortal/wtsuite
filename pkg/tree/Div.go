package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Div struct {
	VisibleTagData
}

func NewDiv(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("div", false, attr, ctx)
	return &Div{visTag}, err
}

func (t *Div) Validate() error {
	return t.ValidateChildren()
}
