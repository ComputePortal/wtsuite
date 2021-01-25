package tree

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
	"github.com/computeportal/wtsuite/pkg/tree/styles"
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
	panic("not available")
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

func (t *LeafTag) CollectClasses(classMap ClassMap) error {
	return nil
}

func (t *LeafTag) CollectStyles(ss styles.DocSheet) error {
	return nil
}

func (t *LeafTag) CollectScripts(idMap IDMap, classMap ClassMap, bundle *scripts.InlineBundle) error {
	return nil
}

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
	if t.parent != nil {
		panic("parent already set")
	}

	t.parent = p
}

func (t *LeafTag) Parent() Tag {
	return t.parent
}

func (t *LeafTag) FinalParent() tokens.FinalTag {
	return t.Parent()
}
