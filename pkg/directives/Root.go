package directives

import (
	"github.com/computeportal/wtsuite/pkg/tree"
	"github.com/computeportal/wtsuite/pkg/tree/scripts"
)

func NewRoot(cache *FileCache, path string, control string, cssUrl string, jsUrl string) (*tree.Root, [][]string, error) {
	_, node, err := BuildFile(cache, path, true, nil)
	if err != nil {
		return nil, nil, err
	}

  return FinalizeRoot(node, control, cssUrl, jsUrl)
}

func FinalizeRoot(node *RootNode, control string, cssUrl string, jsUrl string) (*tree.Root, [][]string, error) {

	root_ := node.tag
	root, ok := root_.(*tree.Root)
	if !ok {
		panic("expected root")
	}

	root.FoldDummy()

	tree.RegisterParents(root)

  if err := root.EvalLazy(); err != nil {
    return nil, nil, err
  }

	idMap := tree.NewIDMap()
	if err := root.CollectIDs(idMap); err != nil {
		return nil, nil, err
	}

	classMap := tree.NewClassMap()
	if err := root.CollectClasses(classMap); err != nil {
		return nil, nil, err
	}

	// bundleableRules is [][]string
	bundleableRules, err := root.CollectStyles(idMap, classMap, cssUrl)
	if err != nil {
		return nil, nil, err
	}

	bundle := scripts.NewInlineBundle()
	if err := root.CollectScripts(idMap, classMap, bundle); err != nil {
		return nil, nil, err
	}

	if control != "" {
		if err := root.ApplyControl(control, jsUrl); err != nil {
			return nil, nil, err
		}
	}

	if err := root.Validate(); err != nil {
		return nil, nil, err
	}

	return root, bundleableRules, nil
}
