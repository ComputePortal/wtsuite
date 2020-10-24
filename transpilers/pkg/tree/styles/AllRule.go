package styles

import (
	tokens "../../tokens/html"
)

// instantiated in Html
type AllRule struct {
	SyncRuleData
}

func NewAllRule(tag StyledTag, attr *tokens.StringDict) (*AllRule, error) {
	rd, err := newSyncRuleData(tag, attr)
	if err != nil {
		return nil, err
	}

	return &AllRule{rd}, nil
}

func (r *AllRule) Synchronize() error {
	return nil
}
