package directives

import (
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

func NewRoot(cache *FileCache, path string, control string, cssUrl string, jsUrl string) (*tree.Root, error) {
	_, node, err := BuildFile(cache, path, true, nil)
	if err != nil {
		return nil, err
	}

  return FinalizeRoot(node, control, cssUrl, jsUrl)
}

func FinalizeRoot(node *RootNode, control string, cssUrl string, jsUrl string) (*tree.Root, error) {

	root_ := node.tag
	root, ok := root_.(*tree.Root)
	if !ok {
		panic("expected root")
	}

	root.FoldDummy()

	tree.RegisterParents(root)

  if err := root.EvalLazy(); err != nil {
    return nil, err
  }

  // checks for uniqueness of id
	idMap := tree.NewIDMap()
	if err := root.CollectIDs(idMap); err != nil {
		return nil, err
	}

  // a similar root.CollectClasses(classMap) used to exist, but it didn't turn out to be very useful

  if cssUrl != "" {
    if err := root.SetStyleURL(cssUrl); err != nil {
      return nil, err
    }
  }

	bundle := scripts.NewInlineBundle()
	if err := root.CollectScripts(bundle); err != nil {
		return nil, err
	}

	if control != "" {
		if err := root.ApplyControl(control, jsUrl); err != nil {
			return nil, err
		}
	}

	if err := root.Validate(); err != nil {
		return nil, err
	}

	return root, nil
}
