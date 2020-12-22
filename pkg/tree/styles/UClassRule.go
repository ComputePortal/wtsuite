package styles

import (
	tokens "../../tokens/html"
)

type UClassRule struct {
	class string // always set a new class
	SyncRuleData
}

func NewUClassRule(tag StyledTag, attr *tokens.StringDict) (*UClassRule, error) {
	rd, err := newSyncRuleData(tag, attr)
	if err != nil {
		return nil, err
	}

	return &UClassRule{"", rd}, nil
}

func (r *UClassRule) Synchronize() error {
	if r.class == "" {
		r.class = newUniqueClass()
		r.tag.SetClasses(append(r.tag.GetClasses(), r.class))
	}

	return nil
}
