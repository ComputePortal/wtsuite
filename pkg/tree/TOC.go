package tree

import (
  "regexp"
  "strconv"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
	"github.com/computeportal/wtsuite/pkg/tree/styles"
)

type TOC struct {
  numbering string
  VisibleTagData
}

func NewTOC(numbering string, attr *tokens.StringDict, ctx context.Context) (Tag, error) {
  prevClass_, ok := attr.Get("class")

  if ok {
    prevClass, err := tokens.AssertString(prevClass_)
    if err != nil {
      return nil, err
    }

    attr.Set("class", tokens.NewValueString(prevClass.Value() + " toc", ctx))
  } else {
    attr.Set("class", tokens.NewValueString("toc", ctx))
  }

  visTag, err := NewVisibleTag("div", false, attr, ctx)
  if err != nil {
    return nil, err
  }

  return &TOC{numbering, visTag}, nil
}

func (t *TOC) Validate() error {
  return nil
}

func (t *TOC) EvalLazy() error {
  // now we need a very special search function
  // we must do a chronological walk from one of the parent nodes
  path := make([]Tag, 0)

  var tag Tag = t
  for tag != nil {
    path = append([]Tag{tag}, path...)
    tag = tag.Parent()
  }

  root := path[0]

  hPattern := regexp.MustCompile(`^h[1-6]$`)

  // now collect all h[1-6] tags in chrono order
  var fnCollect func(tag_ Tag, path_ []Tag) []Tag = nil

  fnCollect = func(tag_ Tag, path_ []Tag) []Tag {
    isRoot := (path_ != nil && len(path_) > 1 && tag_ == path_[0])

    res := make([]Tag, 0)

    children := tag_.Children()

    if isRoot {
      for i, child := range children {
        if isRoot && child == path_[1] {
          children = children[i:]
          break
        }
      }
    }

    for _, child := range children {
      if hPattern.MatchString(child.Name()) {
        res = append(res, child)
      } else {
        if isRoot {
          res = append(res, fnCollect(child, path_[1:])...)
        } else {
          res = append(res, fnCollect(child, nil)...)
        }
      }
    }

    return res
  }

  hTags := fnCollect(root, path)

  ctx := t.Context()


  tocNumbering := NewTOCNumbering()

  for _, hTag := range hTags {
    hNumber_ := hTag.Name()[1:]
    hNumber, err := strconv.ParseInt(hNumber_, 10, 64)
    if err != nil {
      panic(err)
    }

    if t.numbering != "" && t.numbering != "none" {
      tocNumbering.Next(int(hNumber))
      numberPrefixTag := tocNumbering.CreateTag(t.numbering, t.Context())

      if err := hTag.InsertChild(0, numberPrefixTag); err != nil {
        panic(err)
      }
    }

    hID := hTag.GetID()
    if hID == "" {
      hID = styles.NewUniqueID()
      hTag.SetID(hID)
    }

    linkAttr := tokens.NewEmptyStringDict(ctx)
    linkAttr.Set("href", tokens.NewValueString("#" + hID, ctx))
    // XXX: classes can't start with number
    linkAttr.Set("class", tokens.NewValueString("h" + hNumber_, ctx))

    linkTag, err := BuildTag("a", linkAttr, ctx)
    if err != nil {
      panic(err)
    }

    for _, hChild := range hTag.Children() {
      linkTag.AppendChild(Copy(hChild, ctx))
    }

    t.AppendChild(linkTag)
  }

  // only now are is __nchildren__ usable
  if t.attributes != nil {
    if _, err := t.attributes.EvalLazy(t); err != nil {
      return err
    }
  }

  return nil
}
