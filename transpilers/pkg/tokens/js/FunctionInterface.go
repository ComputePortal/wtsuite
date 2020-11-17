package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
	"../patterns"
)

type FunctionInterface struct {
	role prototypes.FunctionRole
	name *VarExpression // can be nil for anonymous functions
	args []*FunctionArgument
	ret  *TypeExpression // can nil for void return ("any" for no return type checking)
}

func NewFunctionInterface(name string, role prototypes.FunctionRole,
	ctx context.Context) *FunctionInterface {
	return &FunctionInterface{
		role,
		NewConstantVarExpression(name, ctx),
		make([]*FunctionArgument, 0),
		nil,
	}
}

func (fi *FunctionInterface) Name() string {
	return fi.name.Name()
}

func (fi *FunctionInterface) Length() int {
	return len(fi.args)
}

func (fi *FunctionInterface) GetVariable() Variable {
	return fi.name.GetVariable()
}

func (fi *FunctionInterface) Context() context.Context {
	return fi.name.Context()
}

func (fi *FunctionInterface) Role() prototypes.FunctionRole {
	return fi.role
}

func (fi *FunctionInterface) SetRole(r prototypes.FunctionRole) {
	fi.role = r
}

func (fi *FunctionInterface) AppendArg(arg *FunctionArgument) {
	fi.args = append(fi.args, arg)
}

// used by parser to gradually fill the interface struct
func (fi *FunctionInterface) SetReturnType(ret *TypeExpression) {
	fi.ret = ret
}

// can be called after resolve names phase
// returns nil if void
// used by return to check type, is used before async (so not a promise)
func (fi *FunctionInterface) GetReturnValue() (values.Value, error) {
  if fi.ret == nil {
    return nil, nil
  } else {
    val, err := fi.ret.EvalExpression()
    if err != nil {
      return nil, err
    }

    return val, nil
  }
}

func (fi *FunctionInterface) Dump() string {
	var b strings.Builder

	// dumping of name can be done here, but writing can't be done below because we need exact control on Function
	if fi.Name() != "" {
		b.WriteString(fi.Name())
	}

	b.WriteString("(")

	for i, arg := range fi.args {
		b.WriteString(arg.Dump(""))

		if i < len(fi.args)-1 {
			b.WriteString(patterns.COMMA)
		}
	}

	b.WriteString(")")

	if fi.ret != nil {
		b.WriteString(fi.ret.Dump(""))
	}

	b.WriteString("\n")

	return b.String()
}

func (fi *FunctionInterface) Write() string {
	var b strings.Builder

	b.WriteString("(")

	for i, arg := range fi.args {
    b.WriteString(arg.Write())

		if i < len(fi.args)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(")")

	return b.String()
}

func (fi *FunctionInterface) performChecks() error {
	// check that arg names are unique, and check that default arguments come last
  detectedDefault := false

	for i, arg := range fi.args {
    if detectedDefault && !arg.HasDefault() {
      errCtx := arg.Context()
      return errCtx.NewError("Error: defaults must come last")
    }

    if arg.HasDefault() {
      detectedDefault = true
    }

    for j, otherArg := range fi.args {
      if i != j {
        if otherArg.Name() == arg.Name() {
          errCtx := context.MergeContexts(otherArg.Context(), arg.Context())
          return errCtx.NewError("Error: argument duplicate name")
        }
      }
    }
	}

  if prototypes.IsGetter(fi) && len(fi.args) != 0 {
    errCtx := fi.args[0].Context()
    return errCtx.NewError("Error: unexpected argument for getter")
  } else if prototypes.IsSetter(fi) && len(fi.args) != 1 {
    errCtx := fi.Context()
    return errCtx.NewError("Error: setter requires exactly one argument")
  }

	return nil
}

func (fi *FunctionInterface) ResolveNames(scope Scope) error {
  if err := fi.performChecks(); err != nil {
    return err
  }

	if fi.ret != nil {
		if err := fi.ret.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	for _, arg := range fi.args {
		if err := arg.ResolveNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) AssertNoDefaults() error {
  for _, arg := range fi.args {
    if err := arg.AssertNoDefault(); err != nil {
      return err
    }
  }

  return nil
}

func (fi *FunctionInterface) GetArgValues() ([]values.Value, error) {
  args := make([]values.Value, len(fi.args))

  for i, fa := range fi.args {
    arg, err := fa.GetValue()
    if err != nil {
      return nil, err
    }

    args[i] = arg
  }

  return args, nil
}

func (fi *FunctionInterface) GetFunctionValue() (*values.Function, error) {
  nOverloads := 1

  for _, arg := range fi.args {
    if arg.HasDefault() {
      nOverloads += 1
    }
  }

  // each argument with a default creates an overload
  argsAndRet := make([][]values.Value, nOverloads)

  retValue, err := fi.ret.EvalExpression()
  if err != nil {
    return nil, err
  }

  if prototypes.IsAsync(fi) {
    if retValue == nil {
      retValue = prototypes.NewVoidPromise(fi.ret.Context())
    } else {
      retValue = prototypes.NewPromise(retValue, fi.ret.Context())
    }
  }

  for i := 0; i < nOverloads; i++ {
    nOverloadArgs := len(fi.args) - (nOverloads - 1 - i)
    argsAndRet[i] = make([]values.Value, nOverloadArgs + 1)

    for j := 0; j < nOverloadArgs; j++ {
      argValue, err := fi.args[j].GetValue()
      if err != nil {
        return nil, err
      }

      argsAndRet[i][j] = argValue
    }

    argsAndRet[i][nOverloadArgs] = retValue
  }
  
  return values.NewOverloadedFunction(argsAndRet, fi.Context()), nil
}

func (fi *FunctionInterface) Eval() error {
  for _, arg := range fi.args {
    if err := arg.Eval(); err != nil {
      return err
    }
	}

	if fi.ret != nil {
		_, err := fi.ret.EvalExpression()
		if err != nil {
			return err
		}
  }

	return nil
}

func (fi *FunctionInterface) UniversalNames(ns Namespace) error {
	for _, arg := range fi.args {
		if err := arg.UniversalNames(ns); err != nil {
			return err
		}
	}

	if fi.ret != nil {
		if err := fi.ret.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) UniqueNames(ns Namespace) error {
	for _, arg := range fi.args {
		if err := arg.UniqueNames(ns); err != nil {
			return err
		}
	}

	if fi.ret != nil {
		if err := fi.ret.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) Walk(fn WalkFunc) error {
  if fi.name != nil {
    if err := fi.name.Walk(fn); err != nil {
      return err
    }
  }

  for _, arg := range fi.args {
    if err := arg.Walk(fn); err != nil {
      return err
    }
  }

  if fi.ret != nil {
    if err := fi.ret.Walk(fn); err != nil {
      return err
    }
  }

  return fn(fi)
}
