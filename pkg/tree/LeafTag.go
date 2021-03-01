package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	//"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

// implements Tag interface
type LeafTag struct {
	parent Tag
	ctx    context.Context
}

func NewLeafTag(ctx context.Context) LeafTag {
	return LeafTag{nil, ctx}
}

func (t *LeafTag) Name() string {
	return ""
}

func (t *LeafTag) GetID() string {
	return ""
}

func (t *LeafTag) SetID(s string) {
	panic("not available")
}

func (t *LeafTag) GetClasses() []string {
	return []string{}
}

func (t *LeafTag) SetClasses(cs []string) {
	panic("not available")
}

func (t *LeafTag) CollectIDs(idMap IDMap) error {
	return nil
}

/*func (t *LeafTag) CollectScripts(bundle *scripts.InlineBundle) error {
	return nil
}*/

func (t *LeafTag) Attributes() *tokens.StringDict {
	return nil
}

func (t *LeafTag) Children() []Tag {
	return []Tag{}
}

func (t *LeafTag) NumChildren() int {
	panic("not available")
}

func (t *LeafTag) AppendChild(child Tag) {
	panic("not available")
}

func (t *LeafTag) InsertChild(i int, child Tag) error {
	panic("not available")
}

func (t *LeafTag) DeleteChild(i int) error {
	panic("not available")
}

func (t *LeafTag) DeleteAllChildren() error {
	panic("not available")
}

func (t *LeafTag) FindID(id *tokens.String) (Tag, int, Tag, bool, error) {
	return nil, 0, nil, false, nil
}

func (t *LeafTag) FoldDummy() {
	return
}

func (t *LeafTag) EvalLazy() error {
  return nil
}

func (t *LeafTag) VerifyElementCount(i int, ecKey string) error {
	return nil
}

func (t *LeafTag) Validate() error {
	// always valid
	return nil
}

func (t *LeafTag) Context() context.Context {
	return t.ctx
}

func (t *LeafTag) Write(indent string, nl, tab string) string {
	panic("not available")
}

func (t *LeafTag) RegisterParent(p Tag) {
  t.parent = p
}

func (t *LeafTag) Parent() Tag {
	return t.parent
}

func (t *LeafTag) Siblings() []Tag {
  parent := t.parent
  if parent == nil {
    panic("parent not yet registered")
  }

  return parent.Children()
}

func (t *LeafTag) LaterSiblings() []Tag {
  parent := t.parent
  if parent == nil {
    panic("parent not yet registered")
  }

  allSiblings := parent.Children()

  for i, s_ := range allSiblings {
    if s, ok := s_.(*LeafTag); ok {
      if s == t {
        return allSiblings[i+1:]
      }
    }
  }

  return []Tag{}
}

func (t *LeafTag) FinalParent() tokens.FinalTag {
	return t.Parent()
}
