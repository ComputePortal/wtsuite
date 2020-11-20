package directives

import (
	"../tree"
	"../tree/scripts"
)

func NewRoot(path string, url string, control string, cssUrl string, jsUrl string) (*tree.Root, [][]string, error) {
	_, node, err := BuildFile(path, "", true)
	if err != nil {
		return nil, nil, err
	}

	root_ := node.tag
	root, ok := root_.(*tree.Root)
	if !ok {
		panic("expected root")
	}

	// postprocessing
	/*if err := root.VerifyElementCount(0, ELEMENT_COUNT); err != nil {
		return nil, nil, nil, err
	}*/

	root.FoldDummy()

	/*if err := root.VerifyElementCount(0, ELEMENT_COUNT_FOLDED); err != nil {
		return nil, nil, nil, err
	}*/

	tree.RegisterParents(root)

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
