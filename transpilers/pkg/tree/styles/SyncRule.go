package styles

import (
	tokens "../../tokens/html"
)

type SyncRule interface {
	Rule
	Synchronize() error // make sure linked tag has correct classes or id
}

type SyncRuleData struct {
	tag StyledTag
	RuleData
}

func newSyncRuleData(tag StyledTag, attr *tokens.StringDict) (SyncRuleData, error) {
	rd, err := newRuleData(attr)
	if err != nil {
		return SyncRuleData{}, err
	}

	return SyncRuleData{tag, rd}, nil
}
