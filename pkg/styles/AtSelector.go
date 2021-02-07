package styles

import (
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree"
)

// at rules can be endlessly deeply nested
type AtSelector struct {
  parent *AtSelector // nil if top level
  key *tokens.String
}

func NewAtSelector(key *tokens.String) *AtSelector {
  return &AtSelector{nil, key}
}

func (s *AtSelector) SetParent(parent *AtSelector) {
  s.parent = parent
}

func collectAtSelectors(sel *AtSelector) []*AtSelector {
  if sel.parent != nil {
    return append(collectAtSelectors(sel.parent), sel)
  } else {
    return []*AtSelector{sel}
  }
}

func (s *AtSelector) Extend(extra *tokens.String) ([]Selector, error) {
  panic("shouldn't be called")
}

func (s *AtSelector) Write() string {
  return s.key.Value()
}

func (s *AtSelector) Match(tag tree.Tag) []tree.Tag {
  // not applicable
  return []tree.Tag{}
}
