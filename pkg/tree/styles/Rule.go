package styles

import (
	tokens "../../tokens/html"
)

type Rule interface {
	Write(indent string) string
	ExpandNested(sel Selector) ([]Rule, error)
}

type RuleData struct {
	attributes *tokens.StringDict
}

func newRuleData(attr *tokens.StringDict) (RuleData, error) {
	return RuleData{attr}, nil
}

func (r *RuleData) Write(indent string) string {
	panic("should've been expanded into a SelectorRule")
}

func (r *RuleData) ExpandNested(sel Selector) ([]Rule, error) {
	leafAttr, rules, err := expandNested(r.attributes, sel)
	if err != nil {
		return nil, err
	}

	if len(leafAttr) > 0 {
		rules = append([]Rule{NewSelectorRule(sel, leafAttr)}, rules...)
	}
	return rules, nil
}
