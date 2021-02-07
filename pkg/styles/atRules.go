package styles

import (
  "strings"

	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type AtRuleGen func(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error)

var _atRules = make(map[string]AtRuleGen)

// list oficial at functions:
// @media
// @charset
// @supports
// @page
// @font-face
// @keyframes

// wtsuite extensions:
// @animation

func registerAtRuleGen(key string, fn AtRuleGen) bool {
  _atRules[key] = fn

  return true
}

// sel == nil is toplevel at rule
func ExpandAtRules(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {
  ctx := key.Context()

  args := strings.Fields(strings.TrimLeft(key.Value(), "@"))

  if len(args) == 0 {
		return nil, ctx.NewError("Error: expected something after the @/.")
  }

  atKey := args[0]

	if fn, ok := _atRules[atKey]; ok {
		return fn(sel, key, attr)
	} else {
		return nil, ctx.NewError("Error: '" + key.Value() + "' at-function not recognized")
	}
}

// generic at rule expansion
func expandGenericAtRule(sel Selector, key *tokens.String, attr *tokens.StringDict, allowLeaf bool) ([]Rule, error) {
  leafAttr, innerRules, err := expandNested(sel, attr)
  if err != nil {
    return nil, err
  }

  if leafAttr.Len() != 0 {
    var leafRule Rule
    if sel == nil {
      if !allowLeaf {
        errCtx := leafAttr.Context()
        return nil, errCtx.NewError("Error: can't have own attributes")
      } else {
        leafRule = NewRule(NewAtSelector(key), leafAttr)
      }
    } else {
      leafRule = NewRule(sel, leafAttr)
    }

    innerRules = append([]Rule{leafRule}, innerRules...)
  }

  // collect all media rules separately
  thisAtSelector := NewAtSelector(key)
  atRules := make([]Rule, 0)
  subRules := make([]Rule, 0)

  for _, innerRule := range innerRules {
    if atRule, ok := innerRule.(*AtRule); ok {
      atRule.SetParent(thisAtSelector)
      atRules = append(atRules, atRule)
    } else {
      subRules = append(subRules, innerRule)
    }
  }

  if len(subRules) > 0 {
    atRules = append(atRules, NewAtRule(thisAtSelector, subRules))
  }

  return atRules, nil
}

func NewMediaRule(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {
  // TODO: check syntax of key
  return expandGenericAtRule(sel, key, attr, false)
}

func NewSupportsRule(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {
  // TODO: check syntax of key
  return expandGenericAtRule(sel, key, attr, false)
}

func NewPageRule(sel Selector, key *tokens.String, attr *tokens.StringDict) ([]Rule, error) {
  return expandGenericAtRule(sel, key, attr, true)
}

var _mediaOk = registerAtRuleGen("media", NewMediaRule)
var _supportsOk = registerAtRuleGen("supports", NewSupportsRule)
var _pageOk = registerAtRuleGen("page", NewPageRule)
