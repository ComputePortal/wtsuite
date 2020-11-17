package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type Interface struct {
	nameExpr *TypeExpression
	parents  []*TypeExpression // can't be nil, can be empty

	members  []*FunctionInterface

  prototypes    []values.Prototype // also used by InstanceOf <interface>

	TokenData
}

func NewInterface(nameExpr *TypeExpression, parents []*TypeExpression,
	ctx context.Context) (*Interface, error) {
  if parents == nil {
    panic("parents can't be nil")
  }

	ci := &Interface{
		nameExpr,
		parents,
		make([]*FunctionInterface, 0),
    make([]values.Prototype, 0),
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
  return []values.Interface{t}, nil
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

func (t *Interface) WriteStatement(indent string) string {
	return ""
}

func (t *Interface) HoistNames(scope Scope) error {
	return nil
}

// uncached check
func (t *Interface) check(other_ values.Interface, ctx context.Context) error {
  v := values.NewInstance(other_, ctx)

  // check each member
  for _, member := range t.members {
    args, err := member.GetArgValues()
    if err != nil {
      return err
    }

    retVal, err := member.GetReturnValue()
    if err != nil {
      return err
    }

    vm, err := v.GetMember(member.Name(), false, ctx)
    if vm == nil && err == nil {
      return ctx.NewError("Error: interface not respected (member " + member.Name() + " not found)")
    }

    if prototypes.IsGetter(member) {
      if retVal == nil {
        panic("getter can't return void, should be checked elsewhere")
      }

      if err := retVal.Check(vm, ctx); err != nil {
        return err
      }
    } else if prototypes.IsSetter(member) {
      if err := vm.SetMember(member.Name(), false, args[0], ctx); err != nil {
        return err
      }
    } else {
      // regular 
      res, err := vm.EvalFunction(args, retVal == nil, ctx)
      if err != nil {
        return err
      }

      if retVal == nil {
        if res != nil {
          return ctx.NewError("Error: interface not respected (member " + member.Name() + " returns non-void, void expected)")
        }
      } else {
        if res == nil {
          return ctx.NewError("Error: interface not respected (member " + member.Name() + " return void, non-void expected)")
        } else {
          if err := retVal.Check(res, ctx); err != nil {
            return err
          }
        }
      }
    } // what about async?
  }

	return nil 
}

// cached Check
func (t *Interface) Check(other_ values.Interface, ctx context.Context) error {
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

		// interface members cant have default arguments, so inner scope is irrelevant
		for _, member := range t.members {
			if err := member.ResolveNames(scope); err != nil {
				return err
			}
		}

		return nil
	}
}

func (t *Interface) EvalStatement() error {
	for _, member := range t.members {
		if err := member.Eval(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Interface) GetInstanceMember(key string, includePrivate bool, ctx context.Context) (values.Value, error) {
  for _, member := range t.members {
    if member.Name() == key {
      if prototypes.IsGetter(member) {
        return member.GetReturnValue()
      } else if prototypes.IsSetter(member) {
        errCtx := ctx
        return nil, errCtx.NewError("Error: is a setter")
      } else {
        return member.GetFunctionValue()
      }
    }
  }

  return nil, nil
}

func (t *Interface) SetInstanceMember(key string, includePrivate bool, arg values.Value, ctx context.Context) error {
  for _, member := range t.members {
    if member.Name() == key {
      if !prototypes.IsSetter(member) {
        return ctx.NewError("Error: not a setter")
      } else {
        args, err := member.GetArgValues()
        if err != nil {
          return err
        }

        return args[0].Check(arg, ctx)
      }
    }
  }

  return ctx.NewError("Error: not a setter")
}

func (t *Interface) ResolveStatementActivity(usage Usage) error {
	return nil
}

func (t *Interface) UniversalStatementNames(ns Namespace) error {
	return nil
}

func (t *Interface) UniqueStatementNames(ns Namespace) error {
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
