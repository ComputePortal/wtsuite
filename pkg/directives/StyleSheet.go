package directives

import (
	"github.com/computeportal/wtsuite/pkg/tree"
)

type StyleSheet interface {
  ApplyExtensions(root *tree.Root) (*tree.Root, error)
}
