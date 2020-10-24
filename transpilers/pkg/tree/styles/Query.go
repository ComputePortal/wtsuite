package styles

import (
	"strings"

	tokens "../../tokens/html"
)

// eg. @supports or @media
type Query struct {
	query    string
	children []Rule
}

func (r *Query) Write(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(r.query)
	b.WriteString("{")
	b.WriteString(NL)

	for _, c := range r.children {
		b.WriteString(c.Write(indent + TAB))
	}

	b.WriteString("}")
	b.WriteString(NL)

	return b.String()
}

func (r *Query) ExpandNested(sel Selector) ([]Rule, error) {
	panic("this is the result of ExpandNested() (can't expand twice)")
}

func newQuery(sel Selector, v tokens.Token, q string) (Query, error) {
	attr, err := tokens.AssertStringDict(v)
	if err != nil {
		return Query{}, err
	}

	leafAttr, rules, err := expandNested(attr, sel)
	if err != nil {
		return Query{}, err
	}

	if len(leafAttr) > 0 {
		rules = append([]Rule{NewSelectorRule(sel, leafAttr)}, rules...)
	}

	return Query{q, rules}, nil
}
