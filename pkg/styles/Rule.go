package styles

import (
  "strings"

	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type Rule interface {
  ExpandNested() ([]Rule, error) // includes self as first
  //ExpandLazy(root *tree.Root) ([]Rule, error) // modifies the root!, adding rules, but keeping all the originals untouched, any remaining lazy values are ignored in the core rule output
  Write(indent string, nl string, tab string) (string, error)
}

type RuleData struct {
  sel Selector
  attr *tokens.StringDict
}

func NewRule(sel Selector, attr *tokens.StringDict) *RuleData {
  return &RuleData{sel, attr}
}

func (r *RuleData) IsLazy() bool {
  return r.attr.ContainsLazy()
}

// so it can be reused by AtRule
func expandNested(sel Selector, attr *tokens.StringDict) (*tokens.StringDict, []Rule, error) {
  ctx := attr.Context()

  rules := make([]Rule, 0)

  leafAttr := tokens.NewEmptyStringDict(ctx)

  // collect leafAttr into leafAttr
  // keys can't be lazy, so we can use regular loop
  if err := attr.Loop(func(key *tokens.String, value_ tokens.Token, last bool) error {
    if tokens.IsNull(value_) {
      // dont do anything
      return nil
    }

    if strings.HasPrefix(key.Value(), "@") {
      value, err := tokens.AssertStringDict(value_)
      if err != nil {
        return err
      }

      // dispatch at rule
      atRuleResult, err := ExpandAtRules(sel, key, value)
      if err != nil {
        return err
      }

      rules = append(rules, atRuleResult...)
    }  else if tokens.IsStringDict(value_) {
      value, err := tokens.AssertStringDict(value_)
      if err != nil {
        return err
      }

      subSels := []Selector{}

      if sel == nil {
        subSels, err = ParseSelectorList(key)
        if err != nil {
          return err
        }
      } else {
        // regular rule extension
        subSels, err = sel.Extend(key)
        if err != nil {
          return err
        }
      }

      for _, subSel := range subSels {
        subRule := NewRule(subSel, value)

        subRules, err := subRule.ExpandNested()
        if err != nil {
          return err
        }

        rules = append(rules, subRules...)
      }
    } else {
      leafAttr.Set(key, value_)
    }

    return nil
  }); err != nil {
    return nil, nil, err
  }

  return leafAttr, rules, nil

}

func (r *RuleData) ExpandNested() ([]Rule, error) {
  leafAttr, rules, err := expandNested(r.sel, r.attr)
  if err != nil {
    return nil, err
  }

  // prepend self if there are some leaf attributes left
  if leafAttr.Len() > 0 {
    leafRule := NewRule(r.sel, leafAttr)

    rules = append([]Rule{leafRule}, rules...)
  }

  return rules, nil
}

func (r *RuleData) writeStart(indent string, nl string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString(r.sel.Write())
	b.WriteString("{")
	b.WriteString(nl)

	return b.String()
}

func (r *RuleData) writeAttributes(indent string, nl string) (string, error) {
	return r.attr.ToString(indent, nl)
}

func (r *RuleData) writeStop(indent string, nl string) string {
	return indent + "}" + nl
}

func (r *RuleData) Write(indent string, nl string, tab string) (string, error) {
	var b strings.Builder

	b.WriteString(r.writeStart(indent, nl))
  inner, err := r.writeAttributes(indent + tab, nl)
  if err != nil {
    return "", err
  }
	b.WriteString(inner)
	b.WriteString(r.writeStop(indent, nl))

	return b.String(), nil
}
