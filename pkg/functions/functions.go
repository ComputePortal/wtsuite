package functions

import (
  "strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	tokens "github.com/computeportal/wtsuite/pkg/tokens/html"
)

// args are evaluated outside PreEval functions, but scope is needed to check Permissive()
type BuiltinFunction func(scope tokens.Scope, args *tokens.Parens, ctx context.Context) (tokens.Token, error)

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
	"merge":             Merge,
	"min":               Min,
	"mix":               Mix,
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
  "reverse":           Reverse,
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
  "find":     Find,
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

func NewUnaryInterface(ctx context.Context) *tokens.Parens {
  return tokens.NewParensInterf([]string{"a"}, nil, ctx)
}

func NewBinaryInterface(ctx context.Context) *tokens.Parens {
  return tokens.NewParensInterf([]string{"a", "b"}, nil, ctx)
}

func NewInterface(names []string, ctx context.Context) *tokens.Parens {
  return tokens.NewParensInterf(names, nil, ctx)
}

// it is up to the caller to have "args" or "interf" be evaluated at this point
func CompleteArgs(args *tokens.Parens, interf *tokens.Parens) ([]tokens.Token, error) {
  if interf == nil {
    for i, alt := range args.Alts() {
      if alt != nil {
        errCtx := args.Values()[i].Context()
        return nil, errCtx.NewError("Error: kwargs not supported")
      }
    }

    return args.Values(), nil
  }

  if args == nil {
    if interf == nil {
      panic("both args and interf can't be 0")
    }

    for i, alt := range interf.Alts() {
      if alt == nil {
        errCtx := interf.Values()[i].Context()
        return nil, errCtx.NewError("Error: doesn't have a default")
      }
    }

    return interf.Alts(), nil
  }

  n := interf.Len()

  if args.Len() > n {
    errCtx := args.Context()
    return nil, errCtx.NewError("Error: expected " + strconv.Itoa(interf.Len()) + " args, got " + strconv.Itoa(args.Len()))
  }

  res := make([]tokens.Token, n)
  for i, _ := range res {
    res[i] = nil
  }

  for i, arg := range args.Values() {
    argAlt := args.Alts()[i]

    if argAlt == nil {
      res[i] = arg
    } else {
      argWord, err := tokens.AssertWord(arg)
      if err != nil {
        return nil, err
      }

      // find the index in interf
      argId := -1
      for j := 0; j < n; j++ {
        interfArg := interf.Values()[j]
        interfArgWord, err := tokens.AssertWord(interfArg)
        if err != nil {
          return nil, err
        }

        if interfArgWord.Value() == argWord.Value() {
          argId = j
          break
        }
      }

      if argId == -1 {
        errCtx := argWord.Context()
        err := errCtx.NewError("Error: kwarg " + argWord.Value() + " not found")
        pos := ""
        for j, interfArgWord_ := range interf.Values() {
          interfArgWord, err := tokens.AssertWord(interfArgWord_)
          if err != nil {
            panic(err)
          }

          pos += interfArgWord.Value() 
          if j < interf.Len() - 1 {
            pos += ", "
          }
        }

        err.AppendString("Info: possibilities are " + pos)

        return nil, err
      }

      res[argId] = argAlt
    }
  }

  // fill remaining with defaults
  for i, r := range res {
    if r != nil {
      continue
    }

    interfAlt := interf.Alts()[i]
    interfWord, err := tokens.AssertWord(interf.Values()[i])
    if err != nil {
      return nil, err
    }

    if interfAlt == nil {
      errCtx := args.Context()
      return nil, errCtx.NewError("Error: arg " + interfWord.Value() + " not specified (doesn't have a default)")
    } else {
      res[i] = interfAlt
    }
  }

  // check the result
  for i, r := range res {
    if r == nil {
      panic("algo error at " + strconv.Itoa(i))
    }
  }

  return res, nil
}

func Eval(scope tokens.Scope, key string, args *tokens.Parens, ctx context.Context) (tokens.Token, error) {
	if fn, ok := preEval[key]; ok {
		evaluated, err := args.EvalAsArgs(scope)
		if err != nil {
			return nil, err
		}

    // if any of the args are lazy, then we can't evaluate the function itself, that would give an error
    if evaluated.AnyLazy() {
      return tokens.NewLazy(func(tag tokens.FinalTag) (tokens.Token, error) {
        evaluated_, err := evaluated.EvalAsArgsLazy(tag)
        if err != nil {
          return nil, err
        }

        return fn(scope, evaluated_, ctx)
      }, ctx), nil
    } else {
      return fn(scope, evaluated, ctx)
    }
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
