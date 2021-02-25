package tree

import (
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

var (
	INLINE           = false
	COMPRESS_NUMBERS = false // always on

	AUTO_LINK = false // optimally convert tags containing 'href' attribute to <a>

	VERBOSITY = 0
)

func buildPermissive(parent Tag, children []*tokens.Tag) error {
  for _, child := range children {
    if child.IsText() {
      childTag := NewText(child.Text(), child.Context())
      parent.AppendChild(childTag)
    } else {
      attr, err := child.Attributes([]string{})
      if err != nil {
        return err
      }

      childTag, err := buildTag(child.Name(), attr, true, child.Context())
      if err != nil {
        return err
      }

      if child.IsScript() {
        childTag.AppendChild(NewText(child.Text(), child.Context()))
      } else {
        if err := buildPermissive(childTag, child.Children()); err != nil {
          return err
        }
      }

      parent.AppendChild(childTag)
    }
  }

  return nil
}

func BuildPermissive(rawTags []*tokens.Tag) (*Root, error) {
  root := NewRoot(rawTags[0].Context())

  if err := buildPermissive(root, rawTags); err != nil {
    return nil, err
  }

  // don't validate, but do register parents

  RegisterParents(root)

  return root, nil
}
