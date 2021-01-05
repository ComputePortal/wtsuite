package styles

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type SelectorRule struct {
	sel  Selector
	attr map[string]string
}

func NewSelectorRule(sel Selector, attr map[string]string) *SelectorRule {
	return &SelectorRule{sel, attr}
}

func (r *SelectorRule) writeStart(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(r.sel.Write())
	b.WriteString("{")
	b.WriteString(patterns.NL)

	return b.String()
}

func (r *SelectorRule) writeAttributes(indent string) string {
	return stringMapToString(r.attr, patterns.NL, indent)
}

func (r *SelectorRule) writeStop(indent string) string {
	return indent + "}" + patterns.NL
}

func (r *SelectorRule) Write(indent string) string {
	var b strings.Builder

	b.WriteString(r.writeStart(indent))
	b.WriteString(r.writeAttributes(indent + patterns.TAB))
	b.WriteString(r.writeStop(indent))

	return b.String()
}

func (r *SelectorRule) ExpandNested(sel Selector) ([]Rule, error) {
	panic("this is the result of ExpandNested() (can't expand twice)")
}
