package tree

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type Generic struct {
	inline bool
	VisibleTagData
}

func NewGeneric(key string, attr *tokens.StringDict, inline bool, ctx context.Context) (Tag, error) {
	visTag, err := NewVisibleTag(key, false, attr, ctx)
	return &Generic{inline, visTag}, err
}

func (t *Generic) Write(indent string, nl, tab string) string {
	if t.inline {
		nl = ""
		tab = ""
		s := t.VisibleTagData.Write("", nl, tab)

		return indent + s
	} else {
		return t.VisibleTagData.Write(indent, nl, tab)
	}
}

func (t *Generic) Validate() error {
	return t.ValidateChildren()
}
