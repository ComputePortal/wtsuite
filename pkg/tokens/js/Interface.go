package js

import (
  "fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"
)

const NewRPCClientMemberName = "newRPCClient"

type Interface struct {
	nameExpr *TypeExpression
	parents  []*VarExpression // can't be nil, can be empty

	members  []*FunctionInterface

  prototypes    []values.Prototype // also used by InstanceOf <interface>

  isRPC bool

	TokenData
}

func NewInterface(nameExpr *TypeExpression, parents []*VarExpression, isRPC bool,
	ctx context.Context) (*Interface, error) {
  if parents == nil {
    panic("parents can't be nil")
  }

	ci := &Interface{
		nameExpr,
		parents,
		make([]*FunctionInterface, 0),
    make([]values.Prototype, 0),
    isRPC,
		TokenData{ctx},
	}

	// change variable so we can register implements classes during resolvenames stage
	ci.nameExpr.variable = NewVariable(ci.Name(), true, ci.nameExpr.Context())
	ci.nameExpr.variable.SetObject(ci)

	return ci, nil
}

func (t *Interface) AddMember(member *FunctionInterface) error {
	// members can have same names, but something must be different (excluding arg names)
	t.members = append(t.members, member)

	return nil
}

func (t *Interface) Name() string {
	return t.nameExpr.Name()
}

func (t *Interface) GetInterfaces() ([]values.Interface, error) {
  interfs := []values.Interface{t}

  for _, parent := range t.parents {
    parentInterfs, err := parent.GetInterface().GetInterfaces()
    if err != nil {
      return nil, err
    }

    // complain if already included
    for _, parentInterf := range parentInterfs {
      for _, interf := range interfs {
        if interf == parentInterf {
          errCtx := parent.Context()
          return nil, errCtx.NewError("Error: interface " + interf.Name() + " extended twice")
        }
      }

      interfs = append(interfs, parentInterf)
    }
  }

  return interfs, nil
}

func (t *Interface) GetPrototypes() ([]values.Prototype, error) {
  return t.prototypes, nil
}

func (t *Interface) GetVariable() Variable {
	return t.nameExpr.GetVariable()
}

func (t *Interface) AddStatement(st Statement) {
	panic("not a block")
}

func (t *Interface) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Interface(")
	b.WriteString(strings.Replace(t.nameExpr.Dump(""), "\n", "", -1))
	b.WriteString(")\n")

  for _, parent := range t.parents {
    b.WriteString(indent + "  extends ")
    b.WriteString(parent.Dump(""))
    b.WriteString("\n")
  }

	for _, member := range t.members {
		if prototypes.IsGetter(member) {
			b.WriteString("\n")
			b.WriteString(indent + "  ")
			b.WriteString("getter ")
		} else if prototypes.IsSetter(member) {
			b.WriteString("\n")
			b.WriteString(indent + "  ")
			b.WriteString("setter ")
		} else {
			b.WriteString(indent + "  ")
		}
		b.WriteString(strings.Replace(member.Dump(), "\n", "", -1))
	}

	return b.String()
}

func (t *Interface) writeRPCNew(indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("static ")
  b.WriteString(NewRPCClientMemberName)
  b.WriteString("(channel,ctx){")
  b.WriteString(nl)

  b.WriteString(indent + tab)
  b.WriteString("let obj={")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  // ctx===ctx_ should be safe enough to see if channel is running on same comms
  b.WriteString("__channel__:function(ctx_){return (ctx===ctx_)?channel:undefined},")
  b.WriteString(nl)

  for _, member := range t.members {
    b.WriteString(member.writeRPCNewEntry(indent + tab + tab, nl, tab))
  }

  b.WriteString(indent + tab)
  b.WriteString("};")
  b.WriteString(nl)

  if len(t.parents) != 0 {
    b.WriteString(indent + tab)
    b.WriteString("Object.assign(obj")
    for _, parent := range t.parents {
      b.WriteString(",")
      b.WriteString(parent.Name())
      b.WriteString(".")
      b.WriteString(NewRPCClientMemberName)
      b.WriteString("(channel,ctx)")
    }
    b.WriteString(");")
  }

  b.WriteString(indent + tab)
  b.WriteString("return obj;")
  b.WriteString(nl)

  b.WriteString(indent)
  b.WriteString("}")
  b.WriteString(nl)

  return b.String()
}

func (t *Interface) writeRPCCallFunctions(indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("static rpcFns(obj,msg,ctx){")
  b.WriteString(nl)

  b.WriteString(indent + tab)
  b.WriteString("let fns={")
  b.WriteString(nl)

  for _, member := range t.members {
    b.WriteString(member.writeRPCCallEntry(indent + tab + tab, nl, tab))
  }

  b.WriteString(indent + tab)
  b.WriteString("};")
  b.WriteString(nl)

  if len(t.parents) != 0 {
    b.WriteString(indent + tab)
    b.WriteString("Object.assign(fns")
    for _, parent := range t.parents {
      b.WriteString(",")
      b.WriteString(parent.Name())
      b.WriteString(".rpcFns(obj,msg,ctx)")
    }
    b.WriteString(");")
  }

  b.WriteString(indent + tab)
  b.WriteString("return fns;")
  b.WriteString(nl)

  b.WriteString(indent)
  b.WriteString("}")
  b.WriteString(nl)

  return b.String()
}

// obj: is the actual instance on which the functions are called
// msg: incoming JSON {id: ..., name: ..., arg0: ..., arg1: ...} (arg0 etc. have already been turned into instances)
// ctx: __rpcContext__
// return value: {value: ...} (if value is left out then void)
func (t *Interface) writeRPCCall(indent string, nl string, tab string) string {
  var b strings.Builder

  b.WriteString(indent)
  b.WriteString("static async rpc(obj,msg,ctx){")
  b.WriteString(nl)

  b.WriteString(indent + tab)
  b.WriteString("try{")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  b.WriteString("let fns=")
  b.WriteString(t.Name())
  b.WriteString(".rpcFns(obj,msg,ctx);")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  b.WriteString("let fn=fns[msg.name];")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  b.WriteString("if(fn==undefined){")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab + tab)
  b.WriteString("throw new Error(msg.name+\" not found\")")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  b.WriteString("}else{await fn()};")
  b.WriteString(nl)

  b.WriteString(indent + tab)
  b.WriteString("}catch(e){")
  b.WriteString(nl)

  b.WriteString(indent + tab + tab)
  b.WriteString("ctx.respond(e);")
  b.WriteString(nl)

  b.WriteString(indent + tab)
  b.WriteString("}")
  b.WriteString(nl)

  b.WriteString(indent)
  b.WriteString("}")
  b.WriteString(nl)

  return b.String()
}

func (t *Interface) WriteStatement(usage Usage, indent string, nl string, tab string) string {
  if !(t.IsRPC() || t.IsUniversal()) {
    fmt.Println(t.IsRPC(), t.IsUniversal())
    return ""
  } else {
    var b strings.Builder

    b.WriteString(indent)
    b.WriteString("class ")
    b.WriteString(t.nameExpr.WriteExpression())
    b.WriteString("{")
    b.WriteString(nl)

    if t.IsUniversal() {
      b.WriteString(indent + tab)
      b.WriteString("static __implementations__=[];")
      b.WriteString(nl)
      // implementations must register themselves, to avoid using classes before they are declared
    }

    if t.IsRPC() {
      b.WriteString(t.writeRPCNew(indent + tab, nl, tab))
      b.WriteString(t.writeRPCCallFunctions(indent + tab, nl, tab))
      b.WriteString(t.writeRPCCall(indent + tab, nl, tab))
    }

    b.WriteString("}")
    b.WriteString(nl)

    return b.String()
  }
}

func (t *Interface) HoistNames(scope Scope) error {
	return nil
}

// v is in instance
// also used for abstract implementations (different error messages)
func checkInterfaceMember(member *FunctionInterface, v values.Value, abstractCheck bool, ctx context.Context) error {
  vm, err := v.GetMember(member.Name(), false, ctx)
  if vm == nil && err == nil {
    if abstractCheck {
      return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
      panic("unexpected")
    } else {
      return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
    }
  } 

  args, err := member.GetArgValues()
  if err != nil {
    return err
  }

  retVal, err := member.GetReturnValue()
  if err != nil {
    return err
  }

  if prototypes.IsGetter(member) {
    if retVal == nil {
      panic("getter can't return void, should be checked elsewhere")
    }

    if vm == nil {
      if abstractCheck {
        return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
        panic("unexpected")
      } else {
        return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
      }
    }

    if err := retVal.Check(vm, ctx); err != nil {
      return err
    }
  } else if prototypes.IsSetter(member) {
    if err := v.SetMember(member.Name(), false, args[0], ctx); err != nil {
      return err
    }
  } else {
    if vm == nil {
      return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
    }

    // regular 
    res, err := vm.EvalFunction(args, retVal == nil, member.Context())
    if err != nil {
      return err
    }

    if retVal == nil {
      if res != nil {
        if abstractCheck {
          return ctx.NewError("Error: abstract member not implemented ("+ member.Name() + " returns non-void, void expected)")
        } else {
          return ctx.NewError("Error: interface not respected (member " + member.Name() + " returns non-void, void expected)")
        }
      }
    } else {
      if res == nil {
        if abstractCheck {
          return ctx.NewError("Error: abstract member not respected (" + member.Name() + " return void, non-void expected)")
        } else {
          return ctx.NewError("Error: interface not respected (member " + member.Name() + " return void, non-void expected)")
        }
      } else {
        if err := retVal.Check(res, ctx); err != nil {
          return err
        }
      }
    }
  } 

  return nil
}

// uncached check
func (t *Interface) check(other_ values.Interface, ctx context.Context) error {
  v := values.NewInstance(other_, ctx)

  // check each member
  for _, member := range t.members {
    if err := checkInterfaceMember(member, v, false, ctx); err != nil {
      return err
    }
  }

  // now check parents
  for _, parent := range t.parents {
    parentInterf := parent.GetInterface()

    // parent interfaces do its caching in turn
    if err := parentInterf.Check(other_, ctx); err != nil {
      return err
    }
  }

	return nil 
}

// cached Check
func (t *Interface) Check(other_ values.Interface, ctx context.Context) error {
  if other, ok := other_.(*Interface); ok {
    if t == other {
      return nil
    }
  }

  if proto, ok := other_.(values.Prototype); ok {
    for _, cached := range t.prototypes {
      if proto == cached {
        return nil
      }
    }

    // first check that proto includes this interface
    protoInterfs, err := proto.GetInterfaces()
    if err != nil {
      return err
    }

    found := false
    for _, protoInterf_ := range protoInterfs {
      if protoInterf, ok := protoInterf_.(*Interface); ok && protoInterf == t {
        found = true
        break
      }
    }

    if !found {
      return ctx.NewError("Error: " + proto.Name() + " doesn't explicitely implement " + t.Name())
    }

    if err = t.check(other_, ctx); err != nil {
      return err
    } else {
      t.prototypes = append(t.prototypes, proto)
      return nil
    }
  } else {
    // should we cache other interface?
    return t.check(other_, ctx)
  }
}

func (t *Interface) ResolveStatementNames(scope Scope) error {
  if t.IsRPC() {
    ActivateMacroHeaders("__checkType__")
  }

	if scope.HasVariable(t.Name()) {
		errCtx := t.Context()
		err := errCtx.NewError("Error: '" + t.Name() + "' already defined " +
			"(interface needs unique name)")
		other, _ := scope.GetVariable(t.Name())
		err.AppendContextString("Info: defined here ", other.Context())
		return err
	} else {
		if err := scope.SetVariable(t.Name(), t.GetVariable()); err != nil {
			return err
		}

    for _, parent := range t.parents {
      if err := parent.ResolveExpressionNames(scope); err != nil {
        return err
      }
    }

		// interface members cant have default arguments
		for _, member := range t.members {
      subScope := NewSubScope(scope)
			if err := member.ResolveNames(subScope); err != nil {
				return err
			}
		}

		return nil
	}
}

func (t *Interface) EvalStatement() error {
  for _, parent := range t.parents {
    parentInterf := parent.GetInterface()
    if parentInterf == nil {
      errCtx := parent.Context()
      return errCtx.NewError("Error: parent is not an interface")
    }

    parentProto := parent.GetPrototype() 
    if parentProto != nil {
      errCtx := parent.Context()
      return errCtx.NewError("Error: parent can't be a prototype")
    }
  }

  // check that members dont exist in parentInterfs
  // getters and setters must be part of the same interface
  interfs, err := t.GetInterfaces() 
  if err != nil {
    return err
  }
  for  _, interf := range interfs {
    if interf == t {
      continue
    }

    for _, member := range t.members {
      if res, err := interf.GetInstanceMember(member.Name(), false, member.Context()); err != nil || res != nil {
        errCtx := member.Context()
        return errCtx.NewError("Error: interface member " + member.Name() + " already exists in extended " + interf.Name())
      }
    }
  }

	for _, member := range t.members {
		if err := member.Eval(); err != nil {
			return err
		}
	}

  if t.isRPC {
    for _, parent := range t.parents {
      parentInterf := parent.GetInterface()

      if !parentInterf.IsRPC() {
        errCtx := parent.Context()
        return errCtx.NewError("Error: rpc interface can't extend non-rpc interface")
      }
    }

    // each return value must be a promise
    for _, member := range t.members {
      if err := member.CheckRPC(); err != nil {
        return err
      }
    }
  }

	return nil
}

func (t *Interface) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  foundGetter := false
  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsGetter(member) {
        foundGetter = true
      }
    }
  }

  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsGetter(member) {
        return member.GetReturnValue()
      } else if prototypes.IsSetter(member) {
        if !foundGetter {
          errCtx := ctx
          return nil, errCtx.NewError("Error: " + t.Name() + "." + key + " is a setter")
        } else {
          continue
        }
      } else {
        return member.GetFunctionValue()
      }
    }
  }

  return nil, nil
}

func (t *Interface) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  foundSetter := false
  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsSetter(member) {
        foundSetter = true
      }
    }
  }

  for _, member := range t.members {
    if member.Name() == key {
      if !prototypes.IsSetter(member) {
        if !foundSetter {
          return ctx.NewError("Error: " + t.Name() + "." + key + " not a setter")
        } else {
          continue
        }
      } else if prototypes.IsSetter(member) {
        args, err := member.GetArgValues()
        if err != nil {
          return err
        }

        return args[0].Check(arg, ctx)
      }
    }
  }

  return ctx.NewError("Error: " + t.Name() + "." + key + " not a setter")
}

// can only be called after eval phase! (because registration is done during eval phase)
func (t *Interface) IsUniversal() bool {
  for _, proto := range t.prototypes {
    if !proto.IsUniversal() {
      return false
    }
  }

  return true
}

func (t *Interface) IsRPC() bool {
  return t.isRPC
}

func (t *Interface) ResolveStatementActivity(usage Usage) error {
	return nil
}

func (t *Interface) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (t *Interface) UniqueStatementNames(ns Namespace) error {
  if t.IsRPC() || t.IsUniversal() {
    if err := ns.ClassName(t.nameExpr.GetVariable()); err != nil {
      return err
    }
  }

	return nil
}

func (t *Interface) Walk(fn WalkFunc) error {
  if err := t.nameExpr.Walk(fn); err != nil {
    return err
  }

  for _, parent := range t.parents {
    if err := parent.Walk(fn); err != nil {
      return err
    }
  }

  for _, member := range t.members {
    if err := member.Walk(fn); err != nil {
      return err
    }
  }

  return fn(t)
}
