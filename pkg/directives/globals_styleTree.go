package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

const PARENT_STYLE = "__pstyle__"

// implements tokens.Token interface
type PStyle struct {
  node Node
  ctx context.Context
}

func setParentStyle(scope Scope, node Node, ctx context.Context) {
  v := functions.Var{&PStyle{node, ctx}, true, true, false, false, ctx}

  scope.SetVar(PARENT_STYLE, v)
}

func (t *PStyle) Eval(scope tokens.Scope) (tokens.Token, error) {
  return t, nil
}

func (t *PStyle) Dump(indent string) string {
  return indent + "__pstyle__"
}

func (t *PStyle) IsSame(other_ tokens.Token) bool {
  if other, ok := other_.(*PStyle); ok {
    return t.node == other.node
  } else {
    return false
  }
}

func (t *PStyle) Context() context.Context {
  return t.ctx
}

func (t *PStyle) Get(scope tokens.Scope, idx_ tokens.Token, ctx context.Context) (tokens.Token, error) {
  if scope.Permissive() && tokens.IsNull(idx_) {
    return tokens.NewNull(ctx), nil
  }

  idx, err := tokens.AssertString(idx_)
  if err != nil {
    return nil, err
  }

  return t.node.SearchStyle(scope, idx, ctx)
}

// search in incoming attr first, the search in parent styles
// used by Math and SVG to search for color
// scope is needed for permissiveness
func SearchStyle(node Node, scope Scope, tagAttr *tokens.StringDict, key string, ctx context.Context) (tokens.Token, error) {
	if styleToken_, ok := tagAttr.Get("style"); ok && !tokens.IsNull(styleToken_) {
		styleToken, err := tokens.AssertStringDict(styleToken_)
		if err != nil {
			return nil, err
		}

		if v, ok := styleToken.Get(key); ok {
			return v, nil
		}
	}

  return node.SearchStyle(scope, tokens.NewValueString(key, ctx), ctx)
}
