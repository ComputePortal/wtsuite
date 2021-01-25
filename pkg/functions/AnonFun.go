package functions

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type AnonFun struct {
	scope       tokens.Scope
	args        *tokens.Parens
	value       tokens.Token
	ctx         context.Context
}

func NewAnonFun(scope tokens.Scope, args *tokens.Parens, value tokens.Token, ctx context.Context) *AnonFun {
	return &AnonFun{scope, args, value, ctx}
}

func (f *AnonFun) Dump(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("AnonFun(")

  b.WriteString(f.args.Dump("  "))
  b.WriteString(")\n")

  return b.String()
}

func (f *AnonFun) Eval(scope tokens.Scope) (tokens.Token, error) {
	return f, nil
}

func (f *AnonFun) EvalLazy(tag tokens.FinalTag) (tokens.Token, error) {
  errCtx := f.Context()
  return nil, errCtx.NewError("Error: function can't be lazily evaluated")
}

func (f *AnonFun) Context() context.Context {
	return f.ctx
}

func (a *AnonFun) IsSame(other tokens.Token) bool {
	if b, ok := other.(*AnonFun); ok {
		if a.args.IsSame(b.args) {
      return true
		}
	}

	return false
}

func CompleteArgsAndFillScope(scope LambdaScope, args_ *tokens.Parens, interf *tokens.Parens) error {
  // whatever comes in has already been oi
  args, err := CompleteArgs(args_, interf)
  if err != nil {
    return err
  }

  fnSet := func(i int, arg tokens.Token) error {
		v := Var{arg, false, true, false, false, arg.Context()}

    argWord, err := tokens.AssertWord(interf.Values()[i])
    if err != nil {
      return err
    }

    if err := scope.SetVar(argWord.Value(), v); err != nil {
      return err
    }

    return nil
  }

  // evaluate the defaults
  for i, arg := range args {
    argDefault_ := interf.Alts()[i]
    if argDefault_ == arg {
      argDefault, err := argDefault_.Eval(scope)
      if err != nil {
        return err
      }

      if err := fnSet(i, argDefault); err != nil {
        return err
      }
    } else {
      if err := fnSet(i, arg); err != nil {
        return err
      }
    }
  }

  return nil
}

func (f *AnonFun) EvalFun(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  var err error
  args_, err = args_.EvalAsArgs(scope)
  if err != nil {
    return nil, err
  }

	lambdaScope := NewLambdaScope(f.scope, scope)

  if err := CompleteArgsAndFillScope(lambdaScope, args_, f.args); err != nil {
    return nil, err
  }

	result, err := f.value.Eval(lambdaScope)
	if err != nil {
		context.AppendContextString(err, "Info: called here", ctx)
		return nil, err
	}

	return result, nil
}

func (f *AnonFun) Len() int {
	return f.args.Len()
}

func IsAnonFun(t tokens.Token) bool {
	_, ok := t.(*AnonFun)
	return ok
}
