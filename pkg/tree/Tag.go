package tree

import (
	"fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"

	//"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

type Tag interface {
	Name() string
	Attributes() *tokens.StringDict

	NumChildren() int
	AppendChild(child Tag)
	InsertChild(i int, child Tag) error // original i is shifted right
	DeleteChild(i int) error
	DeleteAllChildren() error

  EvalLazy() error
	FoldDummy()
	VerifyElementCount(i int, ecKey string) error

	Children() []Tag
	Context() context.Context

	RegisterParent(Tag) // throws error if already set
	Parent() Tag
  FinalParent() tokens.FinalTag // exactly the same as Parent(), but just different interface
  Siblings() []Tag
  LaterSiblings() []Tag // tags that come after!

	//CollectScripts(bundle *scripts.InlineBundle) error

	Validate() error
	Write(indent string, nl, tab string) string
	GetID() string
	SetID(id string)
}

type tagData struct {
	name        string // eg. 'html' or 'head'
	id          string // unique, used for referral, and in js
	selfClosing bool
	attributes  *tokens.StringDict
	parent      Tag
	children    []Tag
	ctx         context.Context
}

func newTag(name string, selfClosing bool, attr *tokens.StringDict, ctx context.Context) (tagData, error) {
	id := ""
	if idToken_, ok := attr.Get("id"); ok {
		if !tokens.IsNull(idToken_) {
			idToken, err := tokens.AssertString(idToken_)
			if err != nil {
				return tagData{}, err
			}

			id = idToken.Value()
		}

		attr.Delete("id")
	}

	if hrefToken, ok := attr.Get("href"); ok {
		if !tokens.IsNull(hrefToken) {
			if _, err := tokens.AssertString(hrefToken); err != nil {
				return tagData{}, err
			}
		} else {
			attr.Delete("href") // is this delete really necessary
		}
	}

	if err := validateAttributes(id, name, attr); err != nil {
		return tagData{}, err
	}

	return tagData{name, id, selfClosing, attr, nil, make([]Tag, 0), ctx}, nil
}

func validateAttributes(id string, name string, attr *tokens.StringDict) error {
	// check attribute values
	return attr.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
		if !tokens.IsBool(value) && !tokens.IsNull(value) && !tokens.IsPrimitive(value) {
			errCtx := value.Context()
			return errCtx.NewError("Error: expected primitive")
		}

		return nil
	})
}

func (t *tagData) Name() string {
	return t.name
}

func (t *tagData) GetID() string {
	return t.id
}

func (t *tagData) SetID(id string) {
	t.id = id
}

// returns ptr, so attributes can be changed in-place
func (t *tagData) Attributes() *tokens.StringDict {
	return t.attributes
}

/*func (t *tagData) CollectScripts(bundle *scripts.InlineBundle) error {
	// scripts themselves dont have children, so can easily override this
	newChildren := make([]Tag, 0)
	for _, child := range t.children {
		if err := child.CollectScripts(bundle); err != nil {
			return err
		}

		switch child.(type) {
		case *Script:
		default:
			newChildren = append(newChildren, child)
		}
	}

	t.children = newChildren

	return nil
}*/

func (t *tagData) collectAttributeClasses() ([]string, error) {
	result := make([]string, 0)

	// get class(es) from explicit class attribute
	var nonUniqueErrorCtx *context.Context = nil
	if classToken, ok := t.attributes.Get("class"); ok {
		errCtx := classToken.Context()
		switch {
		case tokens.IsString(classToken):
			classStrToken, err := tokens.AssertString(classToken)
			if err != nil {
				panic(err)
			}
			classes := strings.Split(classStrToken.Value(), " ")
			result = append(result, classes...)
		case tokens.IsStringList(classToken):
			classLstToken, err := tokens.AssertList(classToken)
			if err != nil {
				panic(err)
			}
			classes, err := classLstToken.GetStrings()
			if err != nil {
				panic(err)
			}

			result = append(result, classes...)
		default:
			return result, errCtx.NewError("Error: invalid class(es)")
		}

		nonUniqueErrorCtx = &errCtx
		t.attributes.Delete("class")
	}

	// check uniqueness
	unique := make(map[string]bool)
	for _, cl := range result {
		if _, ok := unique[cl]; ok {
			errCtx := t.Context()
			if nonUniqueErrorCtx != nil {
				errCtx = *nonUniqueErrorCtx
			}
			return result, errCtx.NewError("Error: non-unique classes")
		} else {
			unique[cl] = true
		}
	}

	return result, nil
}

func writeAttributes(attr *tokens.StringDict) string {
	var b strings.Builder

	if err := attr.Loop(func(key *tokens.String, val tokens.Token, last bool) error {
		// val can also be null, in which case we skip writing it
		if tokens.IsNull(val) || tokens.IsFalseBool(val) {
			return nil
		}

		k := key.Value()
		switch {
    case strings.HasSuffix(k, "!"):
      panic(k + " is an invalid attribute key")
		case k == "style" && !tokens.IsString(val):
      // style can be dict
      value, err := tokens.AssertStringDict(val)
      if err != nil {
        return err
      }

      // should be on single line
      vStr, err := value.ToString("", "")
      if err != nil {
        return err
      }

      b.WriteString(" style=\"")
      b.WriteString(vStr)
      b.WriteString("\"")
		case tokens.IsTrueBool(val):
			b.WriteString(" ")
			b.WriteString(k)
		default:
			value, err := tokens.AssertPrimitive(val)
			if err != nil {
				if tokens.IsList(val) {
          errCtx := val.Context()
          return errCtx.NewError("Error: a list can't be used as a final attribute")
				}
				return err
			}

			v := value.Write()

			b.WriteString(" ")
			b.WriteString(k)
			if v != "" { // empty value is a flag, and doesnt need to be printed
				b.WriteString("=\"")
				b.WriteString(v)
				b.WriteString("\"")
			}
		}

		return nil
	}); err != nil {
		panic("should've been caught in validate\n" + err.Error())
	}

	return b.String()
}

func (t *tagData) writeAttributes() string {
  return writeAttributes(t.attributes)
}

func (t *tagData) writeStartStop(wrapAutoHref bool, indent string, stop bool) string {
	var b strings.Builder

	name := t.name
	hrefToken, hasHref := t.attributes.Get("href")
	if wrapAutoHref && AUTO_LINK && hasHref && !tokens.IsNull(hrefToken) {
		// actually print a
		name = "a"
	}

	b.WriteString(indent)
	b.WriteString("<")
	if stop {
		b.WriteString("/")
	}
	b.WriteString(name)

	if stop {
		b.WriteString(">")
		return b.String()
	}

	b.WriteString(t.writeAttributes())

	if t.selfClosing {
		if patterns.IsCompactSelfClosing(t.name) {
			b.WriteString(">")
		} else {
			b.WriteString("/>")
		}
	} else {
		b.WriteString(">")
	}

	return b.String()
}

func (t *tagData) writeStart(wrapAutoHref bool, indent string) string {
	return t.writeStartStop(wrapAutoHref, indent, false)
}

func (t *tagData) writeStop(wrapAutoHref bool, indent string) string {
	return t.writeStartStop(wrapAutoHref, indent, true)
}

func (t *tagData) write(wrapAutoHref bool, indent string, nl, tab string) string {
	var b strings.Builder

	b.WriteString(t.writeStart(wrapAutoHref, indent))

	if t.selfClosing {
		return b.String()
	} else {
		b.WriteString(nl)

		b.WriteString(t.writeChildren(indent+tab, nl, tab))

		b.WriteString(t.writeStop(wrapAutoHref, indent))
	}

	return b.String()
}

func (t *tagData) Write(indent string, nl, tab string) string {
	return t.write(false, indent, nl, tab)
}

func AssertUniqueID(t Tag, ctx context.Context) (idToken *tokens.String, err error) {
	if vis, ok := t.(VisibleTag); ok {
		if vis.GetID() == "" {
			vis.SetID(NewUniqueID())
		}
		return tokens.NewValueString(vis.GetID(), ctx), nil
	} else {
		if idToken_, ok := t.Attributes().Get("id"); ok {
			idToken, err = tokens.AssertString(idToken_)
			if err != nil {
				return nil, err
			}
		} else {
			idToken = tokens.NewValueString(NewUniqueID(), ctx)

			t.Attributes().Set(tokens.NewValueString("id", ctx), idToken)
		}
	}

	return idToken, nil
}

func (t *tagData) NumChildren() int {
	return len(t.children)
}

func (t *tagData) AppendChild(child Tag) {
  child.RegisterParent(t)

	t.children = append(t.children, child)
}

func (t *tagData) InsertChild(i int, child Tag) error {
  child.RegisterParent(t)

	if i > len(t.children) {
		errCtx := context.MergeContexts(child.Context())
		err := errCtx.NewError("Error: trying to insert child at bad index " + fmt.Sprintf("(i = %d)", i))
		panic(err)
		return err
	}

	children := make([]Tag, 0)

	b := false
	for i_, child_ := range t.children {
		if i_ == i {
			b = true
			children = append(children, child)
		}

		children = append(children, child_)
	}

	if !b {
		children = append(children, child)
	}

	t.children = children

	return nil
}

func (t *tagData) DeleteChild(i int) error {
	if i > len(t.children) {
		errCtx := t.Context()
		return errCtx.NewError("Error: trying to delete child from bad index")
	}

	children := make([]Tag, 0)

	for i_, child_ := range t.children {
		if i_ != i {
			children = append(children, child_)
		}
	}

	t.children = children

	return nil
}

func (t *tagData) DeleteAllChildren() error {
	t.children = []Tag{}
	return nil
}

func (t *tagData) EvalLazy() error {
  if t.attributes != nil {
    if _, err := t.attributes.EvalLazy(t); err != nil {
      return err
    }
  }

  for _, child := range t.children {
    if err := child.EvalLazy(); err != nil {
      return err
    }
  }

  return nil
}

func (t *tagData) FoldDummy() {
	children := make([]Tag, 0)

	for _, child := range t.children {
		child.FoldDummy()
		switch t := child.(type) {
		case *Dummy:
			children = append(children, t.children...)
		default:
			children = append(children, child)
		}
	}

	t.children = children
}

func (t *tagData) VerifyElementCount(i int, ecKey string) error {
	attr := t.Attributes()
	elementCount_, ok := attr.Get(ecKey)
	if !ok {
		errCtx := t.Context()
		return errCtx.NewError("Internal Error: tag doesnt have an " + ecKey)
	}

	elementCount, err := tokens.AssertInt(elementCount_)
	if err != nil {
		return err
	}

	if elementCount.Value() != i {
		errCtx := t.Context()
		return errCtx.NewError(fmt.Sprintf("Internal Error: inconsistent %s, expected %d got %d\n", ecKey, i, elementCount.Value()))

	}

	for i, child := range t.children {
		if err := child.VerifyElementCount(i, ecKey); err != nil {
			return err
		}
	}

	return nil
}

func (t *tagData) writeChildren(indent string, nl, tab string) string {
	var b strings.Builder

	for _, child := range t.children {
		b.WriteString(child.Write(indent, nl, tab))
		b.WriteString(nl)
	}

	return b.String()
}

func (t *tagData) ValidateChildren() error {
	for _, child := range t.children {
		if t.selfClosing {
			errCtx := t.Context()
			return errCtx.NewError("Error: selfclosing cant have children")
		}
		if err := child.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (t *tagData) Children() []Tag {
	return t.children
}

func (t *tagData) Context() context.Context {
	return t.ctx
}

func (t *tagData) RegisterParent(p Tag) {
	t.parent = p
}

func (t *tagData) Validate() error {
  return nil
}

func (t *tagData) Parent() Tag {
	return t.parent
}

func (t *tagData) Siblings() []Tag {
  parent := t.parent
  if parent == nil {
    panic("parent not yet registered")
  }

  return parent.Children()
}

func (t *tagData) LaterSiblings() []Tag {
  parent := t.parent
  if parent == nil {
    panic("parent not yet registered")
  }

  allSiblings := parent.Children()

  for i, s_ := range allSiblings {
    if s, ok := s_.(*tagData); ok {
      if s == t {
        return allSiblings[i+1:]
      }
    }
  }

  return []Tag{}
}

func (t *tagData) FinalParent() tokens.FinalTag {
  return t.Parent()
}

func RegisterParents(root Tag) {
	if _, ok := root.(*LeafTag); ok {
		return
	}

	for _, child := range root.Children() {
		child.RegisterParent(root)

		RegisterParents(child)
	}
}

func FindID(tag Tag, idToken *tokens.String) (Tag, int, Tag, bool, error) {
	id := idToken.Value()

	var parentTag Tag = nil
	var pi int = 0
	var idTag Tag = nil
	var found bool = false

	for i, child := range tag.Children() {
		if child.GetID() == id {
			if found {
				errCtx := context.MergeContexts(idTag.Context(), child.Context())
				return nil, 0, nil, false, errCtx.NewError("Error: id is duplicate")
			}

			parentTag = tag
			pi = i
			idTag = child
			found = true
		}
	}

	for _, child := range tag.Children() {
		resParent, resI, resChild, resFound, resErr := FindID(child, idToken)
		if resErr != nil {
			return nil, 0, nil, false, resErr
		}

		if resFound {
			if found {
				errCtx := context.MergeContexts(idTag.Context(), resChild.Context())
				return nil, 0, nil, false, errCtx.NewError("Error: id is duplicate")
			}

			parentTag = resParent
			pi = resI
			idTag = resChild
			found = true
		}
	}

	return parentTag, pi, idTag, found, nil
}

func WalkText(current Tag, prev []Tag, fn func([]Tag, string) error) error {
	xpath := append(prev, current)
	for _, child_ := range current.Children() {
		switch child := child_.(type) {
		case *Text:
			if err := fn(xpath, child.value); err != nil {
				return err
			}
		default:
			if err := WalkText(child, xpath, fn); err != nil {
				return err
			}
		}
	}

	return nil
}

// superficial copy, used by TOC
func Copy(tag_ Tag, ctx context.Context) Tag {
  switch tag := tag_.(type) {
  case *Text:
    return NewText(tag.Value(), ctx)
  default:
    name := tag_.Name()

    attr_ := tag_.Attributes()

    var attr *tokens.StringDict = nil

    if attr_ != nil {
      attr = tokens.NewEmptyStringDict(ctx)

      // superficial copy of stringdict
      attr_.Loop(func(key *tokens.String, value tokens.Token, last bool) error {
        if key.Value() != "id" {
          attr.Set(key, value)
        }

        return nil
      })
    }

    cpy, err := BuildTag(name, attr, ctx)
    if err != nil {
      panic(err)
    }

    children_ := tag_.Children()

    for _, child := range children_ {
      cpy.AppendChild(Copy(child, ctx))
    }

    return cpy
  }
}
