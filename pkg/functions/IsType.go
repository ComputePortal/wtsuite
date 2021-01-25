package functions

import (
	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

var isTypeTable = map[string]BuiltinFunction{
	"bool":   IsBool,
	"color":  IsColor,
	"dict":   IsDict,
	"float":  IsFloat,
	"int":    IsInt,
	"list":   IsList,
	"null":   IsNull,
	"string": IsString,
}

func IsBool(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsBool(args[0]), ctx), nil
}

func IsColor(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsColor(args[0]), ctx), nil
}

func IsDict(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsDict(args[0]), ctx), nil
}

func IsFloat(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsFloat(args[0]), ctx), nil
}

func IsFunction(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsAnyFunction(args[0]) || IsAnonFun(args[0]), ctx), nil
}

func IsInt(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsInt(args[0]), ctx), nil
}

func IsList(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsList(args[0]), ctx), nil
}

func IsNull(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsNull(args[0]), ctx), nil
}

// everything except null and undefined
func IsVar(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	arg1, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	name, err := tokens.AssertString(arg1)
	if err != nil {
		return nil, err
	}

	fn := tokens.NewFunction("get", []tokens.Token{name, tokens.NewNull(ctx)}, ctx)

	res, err := fn.Eval(scope)
	if err != nil {
		panic(err)
	}

	resIsVar := !tokens.IsNull(res)

	return tokens.NewValueBool(resIsVar, ctx), nil
}

func IsString(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsString(args[0]), ctx), nil
}

func IsType(scope tokens.Scope, args_ *tokens.Parens, ctx context.Context) (tokens.Token, error) {
  args, err := CompleteArgs(args_, nil)
  if err != nil {
    return nil, err
  }

	if len(args) != 2 {
		return nil, ctx.NewError("Error: excepted 2 arguments")
	}

	typeToken, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	if tfn, ok := isTypeTable[typeToken.Value()]; ok {
		return tfn(scope, tokens.NewParens(args[0:1], nil, ctx), ctx)
	} else {
		errCtx := typeToken.Context()
		err := errCtx.NewError("Error: invalid type")
		hint := "Hint, valid types: "
		for k, _ := range isTypeTable {
			hint += k + ", "
		}
		err.AppendString(hint[0 : len(hint)-2])
		return nil, err
	}
}
