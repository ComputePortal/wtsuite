package styles

import (
	tokens "../../tokens/html"
)

type TagRule struct {
	name string // eg. html or body
	SyncRuleData
}

func NewTagRule(tag StyledTag, attr *tokens.StringDict) (*TagRule, error) {
	rd, err := newSyncRuleData(tag, attr)
	if err != nil {
		return nil, err
	}

	return &TagRule{tag.Name(), rd}, nil
}

func (r *TagRule) Synchronize() error {
	return nil
}
