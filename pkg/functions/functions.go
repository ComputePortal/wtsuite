package functions

import (
	"strings"

	"../tokens/context"
	tokens "../tokens/html"
)

// args are evaluated outside PreEval functions, but scope is needed to check Permissive()
type BuiltinFunction func(scope tokens.Scope, args []tokens.Token, ctx context.Context) (tokens.Token, error)

var preEval = map[string]BuiltinFunction{
	"add":               Add,
	"caps":              Caps,
	"ceil":              Ceil,
	"contains":          Contains,
	"cos":               Cos,
	"darken":            Darken,
	"dict":              Dict,
	"dir":               Dir,
	"div":               Div,
	"dump":              Dump,
	"eq":                EQ,
	"error":             Error,
	"float":             Float,
	"floor":             Floor,
	"ge":                GE,
	"get":               Get, // differs from the get(string, [fallback]) function (see directives)
	"gt":                GT,
	"int":               Int,
	"invert":            Invert,
	"isbool":            IsBool,
	"iscolor":           IsColor,
	"isdict":            IsDict,
	"isfloat":           IsFloat,
	"isfunction":        IsFunction,
	"isint":             IsInt,
	"islist":            IsList,
	"isnull":            IsNull,
	"issame":            IsSame,
	"isstring":          IsString,
	"istype":            IsType,
	"items":             Items,
	"join":              Join,
	"keys":              Keys,
	"le":                LE,
	"len":               Len,
	"lighten":           Lighten,
	"list":              List,
	"lower":             Lower,
	"lt":                LT,
	"max":               Max,
	"max-screen-width":  MaxScreenWidth,  // produces string key for media query
	"min-screen-width":  MinScreenWidth,  // produces string key for media query
	"max-screen-height": MaxScreenHeight, // produces string key for media query
	"min-screen-height": MinScreenHeight, // produces string key for media query
	"mix":               Mix,
	"merge":             Merge,
	"min":               Min,
	"mod":               Mod,
	"mul":               Mul,
	"ne":                NE,
	"neg":               Neg,
	"not":               Not,
	"path":              Path,
	"pi":                Pi,
	"pow":               Pow,
	"px":                Px,
	"rad":               Rad, // degrees to rad function
	"rand":              Rand,
	"replace":           Replace,
	"round":             Round,
	"seq":               Seq,
	"sin":               Sin,
	"slice":             Slice,
	"split":             Split,
	"sqrt":              Sqrt,
	"str":               Str,
	"sub":               Sub,
	"svg-path-pos":      SVGPathPos,
	"tan":               Tan,
	"uid":               UniqueID,
	"upper":             Upper,
	"values":            Values,
	"xor":               Xor,
	"year":              Year,
}

var postEval = map[string]BuiltinFunction{
	"and":      And,
	"eval":     EvalFun,
	"filter":   Filter,
	"function": NewFun,
	"ifelse":   IfElse,
	"isvar":    IsVar,
	"map":      Map,
	"or":       Or,
	"sort":     Sort,
}

func HasFun(key string) bool {
	if _, ok := preEval[key]; ok {
		return true
	} else if _, ok := postEval[key]; ok {
		return true
	} else {
		return false
	}
}

func EvalArgs(scope tokens.Scope, args []tokens.Token) ([]tokens.Token, error) {
	evaluated := make([]tokens.Token, len(args))

	for i, arg := range args {
		var err error
		evaluated[i], err = arg.Eval(scope)
		if err != nil {
			return nil, err
		}
	}

	return evaluated, nil
}

func Eval(scope tokens.Scope, key string, args []tokens.Token, ctx context.Context) (tokens.Token, error) {
	if fn, ok := preEval[key]; ok {
		evaluated, err := EvalArgs(scope, args)
		if err != nil {
			return nil, err
		}

		return fn(scope, evaluated, ctx)
	} else if fn, ok := postEval[key]; ok {
		return fn(scope, args, ctx)
	} else {
		err := ctx.NewError("Error: unknown function \"" + key + "\"")
		return nil, err
	}
}

func ListValidNames() string {
	var b strings.Builder

	b.WriteString("\n")

	for k, _ := range preEval {
		b.WriteString(k)
		b.WriteString(" (built-in function)")
		b.WriteString("\n")
	}

	b.WriteString("\n")

	for k, _ := range postEval {
		b.WriteString(k)
		b.WriteString(" (built-in function)")
		b.WriteString("\n")
	}

	return b.String()
}
