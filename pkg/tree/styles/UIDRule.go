package styles

import (
	tokens "../../tokens/html"
)

type UIDRule struct {
	SyncRuleData
}

func NewUIDRule(tag StyledTag, attr *tokens.StringDict) (*UIDRule, error) {
	rd, err := newSyncRuleData(tag, attr)
	if err != nil {
		return nil, err
	}

	return &UIDRule{rd}, nil
}

func (r *UIDRule) ToUClassRule() *UClassRule {
	return &UClassRule{"", SyncRuleData{r.tag, RuleData{r.attributes}}}
}

func (r *UIDRule) Synchronize() error {
	if r.tag.GetID() == "" {
		r.tag.SetID(NewUniqueID())
	}

	return nil
}
