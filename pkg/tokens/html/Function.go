package html

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Function struct {
	name string
	args *Parens
	TokenData
}

func NewValueFunction(name string, args *Parens, ctx context.Context) *Function {
	return &Function{name, args, TokenData{ctx}}
}

func NewFunction(name string, args []Token, ctx context.Context) *Function {
	return NewValueFunction(name, NewParens(args, nil, ctx), ctx)
}

func (t *Function) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent + "Function " + t.name + "\n")

	for i, arg := range t.args.values {
    alt := t.args.alts[i]

    if alt == nil {
      b.WriteString(arg.Dump(indent + "  "))
    } else {
      b.WriteString(arg.Dump(indent + "  "))
      b.WriteString(alt.Dump(indent + "= "))
    }
	}

	return b.String()
}

func IsAnyFunction(t Token) bool {
	_, ok := t.(*Function)
	return ok
}

func IsFunction(t Token, name string) bool {
	if fn, ok := t.(*Function); ok {
		return fn.name == name
	}

	return false
}

func AssertFunction(t Token) (*Function, error) {
	if fn, ok := t.(*Function); ok {
		return fn, nil
	} else {
		errCtx := t.Context()
		return nil, errCtx.NewError("Error: expected a function")
	}
}

func (t *Function) Eval(scope Scope) (Token, error) {
  res, err := scope.Eval(t.name, t.args, t.Context())
  if err != nil {
    return nil, err
  }

  if _, ok := res.(*Function); ok {
    errCtx := res.Context()
    err := errCtx.NewError("Internal Error: result of an eval can't be another function")
    panic(err)
  }

  return res, nil
}

func (t *Function) EvalLazy(tag FinalTag) (Token, error) {
  errCtx := t.Context()
  return nil, errCtx.NewError("Error: unable to evaluate lazily")
}

func (a *Function) IsSame(other Token) bool {
	if b, ok := other.(*Function); ok {
		if a.name == b.name {
      if a.args.IsSame(b.args) {
        return true
      }
		}
	}

	return false
}

func (t *Function) Args() *Parens {
	return t.args
}
