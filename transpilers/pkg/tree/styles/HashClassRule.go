package styles

import (
	tokens "../../tokens/html"
	"../../tokens/raw"
)

type HashClassRule struct {
	hash string
	SyncRuleData
}

func NewHashClassRule(tag StyledTag, attr *tokens.StringDict) (*HashClassRule, error) {
	hash := raw.ShortHash(attr.Dump(""))

	rd, err := newSyncRuleData(tag, attr)
	if err != nil {
		return nil, err
	}

	return &HashClassRule{hash, rd}, nil
}

func (r *HashClassRule) ToUClassRule() *UClassRule {
	return &UClassRule{r.hash, SyncRuleData{r.tag, RuleData{r.attributes}}}
}

func (r *HashClassRule) Synchronize() error {
	r.tag.SetClasses(append(r.tag.GetClasses(), r.hash))

	return nil
}
