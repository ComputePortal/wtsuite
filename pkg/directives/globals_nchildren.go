package directives

import (
	"github.com/computeportal/wtsuite/pkg/functions"
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

const NCHILDREN = "__nchildren__"
const NSIBLINGS = "__nsiblings__"

// TODO: the tag isn't necessarily the right one
func setNChildren(scope Scope, ctx context.Context) {
  val := tokens.NewLazy(func(tag tokens.FinalTag) (tokens.Token, error) {
    return tokens.NewValueInt(tag.NumChildren(), ctx), nil
  }, ctx)

  valVar := functions.Var{val, true, true, false, false, ctx}

  if err := scope.SetVar(NCHILDREN, valVar); err != nil {
    panic(err)
  }
}

// including self
func setNSiblings(scope Scope, ctx context.Context) {
  val := tokens.NewLazy(func(tag tokens.FinalTag) (tokens.Token, error) {
    parent := tag.FinalParent()
    if parent == nil {
      return tokens.NewValueInt(1, ctx), nil
    } else {
      return tokens.NewValueInt(parent.NumChildren(), ctx), nil
    }
  }, ctx)

  valVar := functions.Var{val, true, true, false, false, ctx}

  if err := scope.SetVar(NSIBLINGS, valVar); err != nil {
    panic(err)
  }
}

func setLazyTagVars(scope Scope, ctx context.Context) {
  setNChildren(scope, ctx)

  setNSiblings(scope, ctx)
}
