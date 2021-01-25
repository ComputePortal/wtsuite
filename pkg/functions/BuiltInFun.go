package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

type BuiltInFun struct {
	name string
	ctx  context.Context
}

func NewBuiltInFun(name string, ctx context.Context) *BuiltInFun {
	return &BuiltInFun{name, ctx}
}

func (f *BuiltInFun) Dump(indent string) string {
	return indent + "BuiltInFun(" + f.name + ")\n"
}

func (f *BuiltInFun) Eval(scope tokens.Scope) (tokens.Token, error) {
	return f, nil
}

func (f *BuiltInFun) EvalLazy(tag tokens.FinalTag) (tokens.Token, error) {
  errCtx := f.Context()
  return nil, errCtx.NewError("Error: function can't be lazily evaluated")
}

func (f *BuiltInFun) Context() context.Context {
	return f.ctx
}

func (a *BuiltInFun) IsSame(other tokens.Token) bool {
	if b, ok := other.(*BuiltInFun); ok {
		return a.name == b.name
	}
	return false
}

func (f *BuiltInFun) EvalFun(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	result, err := scope.Eval(f.name, args, ctx)
	if err != nil {
		context.AppendContextString(err, "Info: function defined here", f.Context())
		return nil, err
	}

	return result, nil
}

func (f *BuiltInFun) Len() int {
	switch f.name {
	case "uid":
		return 0
	case "int", "isbool", "iscolor", "isdict", "isfloat", "isint", "islist", "isnull", "isstring", "items", "keys", "not", "len", "str", "values", "element", "constructor", "sqrt", "caps":
		return 1
	case "add", "and", "div", "eq", "ge", "gt", "istype", "le", "lt", "merge", "mul", "ne", "or", "issame", "sub", "xor", "sort", "map", "filter":
		return 2
	case "ifelse":
		return 3
	case "dict":
		return -1
	default:
		return -1
	}
}
