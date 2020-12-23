package js

import (
  "strconv"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
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
func (fi *FunctionInterface) getReturnValue() (values.Value, error) {
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

// post async
func (fi *FunctionInterface) GetReturnValue() (values.Value, error) {
  ret, err := fi.getReturnValue()
  if err != nil {
    return nil, err
  }

  if prototypes.IsAsync(fi) {
    if ret == nil {
      return prototypes.NewVoidPromise(fi.Context()), nil
    } else {
      return prototypes.NewPromise(ret, fi.ret.Context()), nil
    }
  } else {
    return ret, nil
  }
}

func (fi *FunctionInterface) IsVoid() bool {
  return fi.ret == nil
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

func (fi *FunctionInterface) writeRPCNewEntry(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(fi.Name())
  b.WriteString(":")
  b.WriteString("async function(")
  for i, _ := range fi.args {
    b.WriteString("arg")
    b.WriteString(strconv.Itoa(i))
    if i < len(fi.args) - 1 {
      b.WriteString(",")
    }
  }
  b.WriteString("){")
  b.WriteString(NL)

  // create the msg
  b.WriteString(indent + TAB)
  b.WriteString("let packet={channel:channel,name:\"")
  b.WriteString(fi.Name())
  b.WriteString("\"")

  for i, arg := range fi.args {
    argName := "arg" + strconv.Itoa(i)
    b.WriteString(",")
    b.WriteString(argName)
    b.WriteString(":")

    te := arg.typeExpr

    interf := te.GetInterface()
    if interf == nil {
      panic("unexpected")
    }

    if interf.IsRPC() {
      // arg could be an already open channel
      b.WriteString("((")
      b.WriteString(argName)
      b.WriteString(".__channel__!==undefined&&")
      b.WriteString(argName)
      b.WriteString(".__channel__(ctx)!==undefined)?")
      b.WriteString("{channel:")
      b.WriteString(argName)
      b.WriteString(".__channel__(ctx)}")
      b.WriteString(":ctx.register(")
      b.WriteString(interf.Name())
      b.WriteString(",")
      b.WriteString(argName)
      b.WriteString("))")
    } else {
      b.WriteString(argName)
    }
  }
  b.WriteString("};")
  b.WriteString(NL)

  // prepare a return value
  promise, err := fi.GetReturnValue()
  if err != nil {
    panic("should've been detected before")
  }

  ret, err := prototypes.GetPromiseContent(promise)
  if err != nil {
    panic("should've been detected before")
  }

  // perform the actual request
  b.WriteString(indent + TAB)
  if ret != nil {
    b.WriteString("let result=")
  }
  b.WriteString("await ctx.request(packet,[")
  // we also need to tell ctx which new channels have been created (so they can be deleted when the request is fullfilled)
  for i, arg := range fi.args {
    te := arg.typeExpr
    interf := te.GetInterface()
    if interf.IsRPC() {
      b.WriteString("packet.arg")
      b.WriteString(strconv.Itoa(i))
      b.WriteString(",")
    }
  }
  b.WriteString("]);")
  b.WriteString(NL)

  // type check?
  if ret != nil {
    b.WriteString(indent + TAB)
    b.WriteString("return ")

    retInterf := values.GetInterface(ret)
    if retInterf == nil {
      panic("unexpected")
    } else if retInterf.IsRPC() {
      b.WriteString(retInterf.Name())
      b.WriteString(".")
      b.WriteString(NewRPCClientMemberName)
      b.WriteString("(parseInt(__checkType__(result.value,Number)),ctx);")
    } else if retInterf.IsUniversal() {
      b.WriteString("__checkType__(result.value,")
      b.WriteString(fi.ret.WriteUniversalRuntimeType())
      b.WriteString(");")
    } else {
      panic("unexpected")
    }

    b.WriteString(NL)
  }

  b.WriteString(indent)
  b.WriteString("},")
  b.WriteString(NL)

  return b.String()
}

func (fi *FunctionInterface) writeRPCCallEntry(indent string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString(fi.Name())
  b.WriteString(":async function(){")
  b.WriteString(NL)

  // each argument is extracted from the msg, and type checked
  for i, arg := range fi.args {
    argName := "arg" + strconv.Itoa(i)

    b.WriteString(indent + TAB)
    b.WriteString("let ") 
    b.WriteString(argName)
    b.WriteString(";")
    b.WriteString(NL)

    te := arg.typeExpr

    interf := te.GetInterface() 
    if interf == nil {
      panic("unexpected")
    } else if interf.IsRPC() {
      b.WriteString(indent + TAB)
      b.WriteString("if(msg.")
      b.WriteString(argName)
      b.WriteString("!==undefined&&msg.")
      b.WriteString(argName)
      b.WriteString(".channel!==undefined){")
      b.WriteString(NL)

      b.WriteString(indent + TAB + TAB)
      b.WriteString("let pair=ctx.channels[parseInt(__checkType__(msg.")
      b.WriteString(argName)
      b.WriteString(".channel,Number))];")
      b.WriteString(NL)
      b.WriteString(indent + TAB + TAB)
      b.WriteString("if(pair==undefined){throw new Error('channel not found')}")
      b.WriteString(NL)
      b.WriteString(indent + TAB + TAB)
      b.WriteString(argName)
      b.WriteString("=pair[1];")
      b.WriteString(NL)

      b.WriteString(indent + TAB)
      b.WriteString("}else{")
      b.WriteString(NL)

      b.WriteString(indent + TAB + TAB)
      b.WriteString(argName)
      b.WriteString("=")
      b.WriteString(interf.Name())
      b.WriteString(".")
      b.WriteString(NewRPCClientMemberName)
      b.WriteString("(parseInt(__checkType__(msg.")
      b.WriteString(argName)
      b.WriteString(",Number)),ctx)")
      b.WriteString(NL)

      b.WriteString(indent + TAB)
      b.WriteString("}")
    } else {
      b.WriteString(indent + TAB);
      b.WriteString(argName)
      b.WriteString("=__checkType__(msg.")
      b.WriteString(argName)
      b.WriteString(",")
      b.WriteString(te.WriteUniversalRuntimeType())
      b.WriteString(");")
    }

    b.WriteString(NL)
  }

  // prepare a return value
  promise, err := fi.GetReturnValue()
  if err != nil {
    panic("should've been detected before")
  }

  ret, err := prototypes.GetPromiseContent(promise)
  if err != nil {
    panic("should've been detected before")
  }

  // now call the actual function
  b.WriteString(indent + TAB)
  if ret != nil {
    b.WriteString("let r=")
  }
  b.WriteString("await obj.")
  b.WriteString(fi.Name())
  b.WriteString("(")
  for i, _ := range fi.args {
    b.WriteString("arg")
    b.WriteString(strconv.Itoa(i))
    if i < len(fi.args) - 1 {
      b.WriteString(",")
    }
  }
  b.WriteString(");")
  b.WriteString(NL)

  b.WriteString(indent + TAB)
  b.WriteString("ctx.respond({id:msg.id,")
  if ret != nil {
    retInterf := values.GetInterface(ret)
    if retInterf == nil {
      panic("unexpected")
    } else if retInterf.IsRPC() {
      b.WriteString("value:ctx.register(")
      b.WriteString(retInterf.Name())
      b.WriteString(",r)")
    } else if retInterf.IsUniversal() {
      b.WriteString("value:r")
    } else {
      panic("unexpected")
    }
  }
  b.WriteString("});")
  b.WriteString(NL)

  b.WriteString(indent)
  b.WriteString("},")
  b.WriteString(NL)

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

  retValue, err := fi.GetReturnValue()
  if err != nil {
    return nil, err
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

func (fi *FunctionInterface) CheckRPC() error {
  ctx := fi.Context()
  // disallow getters and setters
  if prototypes.IsSetter(fi) {
    return ctx.NewError("Error: rpc member can't be setter")
  }

  if prototypes.IsGetter(fi) {
    return ctx.NewError("Error: rpc member can't be getter")
  }

  retVal, err := fi.GetReturnValue()
  if err != nil {
    return err
  }

  if retVal == nil {
    return ctx.NewError("Error: rpc member return value expected Promise, got void")
  }

  if !prototypes.IsPromise(retVal) {
    return ctx.NewError("Error: rpc member return value expected Promise, got " + retVal.TypeName())
  }

  promiseContent, err := prototypes.GetPromiseContent(retVal)
  if err != nil {
    return err
  }

  if promiseContent != nil {
    promiseContentInterf := values.GetInterface(promiseContent)
    if promiseContentInterf == nil {
      errCtx := retVal.Context()
      return errCtx.NewError("Error: rpc member returns non-universal/non-rpc Promise")
    } else if !(promiseContentInterf.IsRPC() || promiseContentInterf.IsUniversal()) {
      errCtx := retVal.Context()
      return errCtx.NewError("Error: rpc member returns non-universal/non-rpc Promise<" + promiseContentInterf.Name() + ">")
    }
  }

  // each arg's interface must be universal or rpc, any is not allowed
  for _, arg := range fi.args {
    argVal, err := arg.GetValue()
    if err != nil {
      return err
    }

    interf := values.GetInterface(argVal)
    if interf == nil {
      errCtx := argVal.Context()
      return errCtx.NewError("Error: expected universal/rpc arg, got " + argVal.TypeName())
    } else if !(interf.IsUniversal() || interf.IsRPC()) {
      errCtx := argVal.Context()
      return errCtx.NewError("Error: expected universal/rpc arg, got " + argVal.TypeName())
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
