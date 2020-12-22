package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Body struct {
	VisibleTagData
}

func NewBody(attr *tokens.StringDict, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag("body", false, attr, ctx)
	return &Body{visTag}, err
}

func (t *Body) Validate() error {
	return t.ValidateChildren()
}

func (t *Body) CollectIDs(idMap IDMap) error {
	if t.id == "" {
		idMap.Set("body", t)
	}

	return t.VisibleTagData.CollectIDs(idMap)
}
