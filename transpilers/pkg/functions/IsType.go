package functions

import (
	"../tokens/context"
	tokens "../tokens/html"
)

type IsTypeFunction func(args []tokens.Token, ctx context.Context) (tokens.Token, error)

var isTypeTable = map[string]IsTypeFunction{
	"bool":   IsBool,
	"color":  IsColor,
	"dict":   IsDict,
	"float":  IsFloat,
	"int":    IsInt,
	"list":   IsList,
	"null":   IsNull,
	"string": IsString,
}

func IsBool(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsBool(args[0]), ctx), nil
}

func IsColor(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsColor(args[0]), ctx), nil
}

func IsDict(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsDict(args[0]), ctx), nil
}

func IsFloat(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsFloat(args[0]), ctx), nil
}

func IsFunction(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsAnyFunction(args[0]) || IsAnonFun(args[0]), ctx), nil
}

func IsInt(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsInt(args[0]), ctx), nil
}

func IsList(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsList(args[0]), ctx), nil
}

func IsNull(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsNull(args[0]), ctx), nil
}

// everything except null and undefined
func IsVar(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	if str, ok := args[0].(*tokens.String); ok && str.WasWord() {
		args[0] = tokens.NewValueString(str.Value(), str.Context())
	}

	arg1, err := args[0].Eval(scope)
	if err != nil {
		return nil, err
	}

	name, err := tokens.AssertString(arg1)
	if err != nil {
		return nil, err
	}

	fn := tokens.NewValueFunction("get", []tokens.Token{name, tokens.NewNull(ctx)}, ctx)

	res, err := fn.Eval(scope)
	if err != nil {
		panic(err)
	}

	resIsVar := !tokens.IsNull(res)

	return tokens.NewValueBool(resIsVar, ctx), nil
}

func IsString(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return tokens.NewValueBool(tokens.IsString(args[0]), ctx), nil
}

func IsType(args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: excepted 2 arguments")
	}

	typeToken, err := tokens.AssertString(args[1])
	if err != nil {
		return nil, err
	}

	if tfn, ok := isTypeTable[typeToken.Value()]; ok {
		return tfn(args[0:1], ctx)
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
