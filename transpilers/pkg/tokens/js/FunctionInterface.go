package js

import (
	"fmt"
	"reflect"
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

func (fi *FunctionInterface) Rest() bool {
	if len(fi.args) == 0 {
		return false
	} else {
		return fi.args[len(fi.args)-1].Rest()
	}
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

func (fi *FunctionInterface) SetReturnType(ret *TypeExpression) {
	fi.ret = ret
}

func (fi *FunctionInterface) Check() error {
	argNameDone := make(map[string]*FunctionArgument)

	// check that arg names are unique
	for _, arg := range fi.args {
		if other, ok := argNameDone[arg.Name()]; ok {
			errCtx := context.MergeContexts(other.Context(), arg.Context())
			return errCtx.NewError("Error: argument duplicate name")
		}

		argNameDone[arg.Name()] = arg
	}

	if fi.Rest() && fi.args[len(fi.args)-1].def != nil {
		errCtx := fi.args[len(fi.args)-1].def.Context()
		return errCtx.NewError("Error: default argument for rest argument no allowed")
	}

	return nil
}

func (fi *FunctionInterface) HasArg(name string) bool {
	for _, arg := range fi.args {
		if arg.Name() == name {
			return true
		}
	}

	return false
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
		if fi.Rest() && i == len(fi.args)-1 {
			b.WriteString(arg.Write(true))
		} else {
			b.WriteString(arg.Write(false))
		}

		if i < len(fi.args)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(")")

	return b.String()
}

func (fi *FunctionInterface) ResolveInterfaceNames(scope Scope) error {
	for _, arg := range fi.args {
		if err := arg.ResolveInterfaceNames(scope); err != nil {
			return err
		}
	}

	if fi.ret != nil {
		if err := fi.ret.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) ResolveNames(outer Scope, inner Scope) error {
	for _, arg := range fi.args {
		if err := arg.ResolveNames(outer, inner); err != nil {
			return err
		}
	}

	if fi.ret != nil {
		if err := fi.ret.ResolveExpressionNames(outer); err != nil {
			return err
		}
	}

	return nil
}

func (fi *FunctionInterface) GenerateArgInstances(stack values.Stack,
	ctx context.Context) ([]values.Value, error) {
	if fi.Rest() {
		return nil, ctx.NewError("Error: cannot generate arg instances with rest")
	}

	res := make([]values.Value, 0)

	for _, arg := range fi.args {
		argInstance, err := arg.GenerateArgInstance(stack, ctx)
		if err != nil {
			return nil, err
		}

		res = append(res, argInstance)
	}

	return res, nil
}

func (fi *FunctionInterface) GenerateReturnInstance(stack values.Stack,
	ctx context.Context) (values.Value, error) {
	if fi.ret == nil {
		errCtx := fi.Context()
		return nil, errCtx.NewError("Error: doesn't have a return value")
	}

	return fi.ret.GenerateInstance(stack, ctx)
}

func (fi *FunctionInterface) EvalInterface(stack values.Stack) error {
	for _, arg := range fi.args {
		if err := arg.EvalInterface(stack); err != nil {
			return err
		}
	}

	if fi.ret != nil {
		retClassVal, err := fi.ret.EvalExpression(stack)
		if err != nil {
			return err
		}

		_, ok := retClassVal.GetClassInterface()
		if !ok {
			errCtx := fi.ret.Context()
			return errCtx.NewError("Error: not a class or interface")
		}
	}

	return nil
}

func (fi *FunctionInterface) EvalArgs(stack values.Stack, args []values.Value, ctx context.Context) error {
	i := 0
	for ; i < len(args); i++ {
		if !fi.Rest() && i > len(fi.args)-1 {
			return ctx.NewError("Error: too many arguments, expected " +
				fmt.Sprintf("%d but got %d", len(fi.args), len(args)))
		} else if fi.Rest() && i == len(fi.args)-1 {
			if err := fi.args[i].EvalRest(stack, args[i:], ctx); err != nil {
				return err
			}

			break
		} else {
			if err := fi.args[i].EvalArg(stack, args[i], ctx); err != nil {
				return err
			}
		}
	}

	// any remaining args must have defaults
	for ; i < len(fi.args); i++ {
		if err := fi.args[i].EvalDef(stack, ctx); err != nil {
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

func (fi *FunctionInterface) ConstrainReturnValue(stack values.Stack, val values.Value, ctx context.Context) (values.Value, error) {
	if fi.ret == nil {
		if val != nil && !values.IsVoid(val) {
			// function and the value context should be in the same file, so merging the contexts shouldn't be a problem
			errCtx := context.MergeContexts(fi.name.Context(), val.Context())
			return nil, errCtx.NewError("Error: no return type specified")
		} else {
			return nil, nil
		}
	} else if val == nil {
		errCtx := fi.name.Context()
		return nil, errCtx.NewError("Error: expected a return value, got void")
	} else {
		return fi.ret.Constrain(stack, val)
	}
}

func (fi *FunctionInterface) IsImplementedBy(proto_ values.Prototype) (string, bool) {

	switch proto := proto_.(type) {
	case *Class:
		return fi.isImplementedByClass(proto)
	case *Enum:
		return fi.IsImplementedBy(proto.cachedExtends)
	case prototypes.BuiltinPrototypeInterface:
		return fi.isImplementedByBuiltinPrototype(proto)
	default:
		panic("not implemented " + reflect.TypeOf(proto_).String())
	}
}

func (fi *FunctionInterface) isImplementedByBuiltinPrototype(proto prototypes.BuiltinPrototypeInterface) (string, bool) {
	builtinFn := proto.FindMember(fi.Name())

	if builtinFn == nil {
		return fi.Name() + " missing", false
	} else {
		return fi.isImplementedByBuiltinFunction(builtinFn)
	}
}

func (fi *FunctionInterface) isImplementedByClass(cl *Class) (string, bool) {
	clMember := cl.FindMember(fi.Name(), false, prototypes.IsSetter(fi))
	if clMember == nil {
		if cl.cachedExtends != nil {
			return fi.IsImplementedBy(cl.cachedExtends)
		} else {
			return fi.Name() + " missing", false
		}
	}

	return fi.isImplementedByOtherInterface(clMember.function.Interface())
}

func (fi *FunctionInterface) isImplementedByOtherInterface(other *FunctionInterface) (string, bool) {
	if fi.Name() != other.Name() {
		panic("should've been caught before")
	}

	if fi.role != other.role {
		return fi.Name() + " differs", false
	}

	if len(fi.args) > len(other.args) {
		return fi.Name() + " differs", false
	} else if len(fi.args) < len(other.args) {
		// remaining other needs default args
		for i := len(fi.args); i < len(other.args); i++ {
			otherArg := other.args[i]
			if otherArg.def == nil {
				return fi.Name() + " differs", false
			}
		}
	}

	for i, arg := range fi.args {
		// if this arg doesnt doesnt have a constraint, neither can the other
		if !arg.IsImplementedByOtherArg(other.args[i]) {
			return fi.Name() + " differs", false
		}
	}

	if fi.ret == nil {
		if other.ret != nil {
			return fi.Name() + "'s return value differs", false
		}
	} else if other.ret == nil {
		return fi.Name() + "'s return value differs", false
	} else {
		if fi.ret.Dump("") != other.ret.Dump("") {
			return fi.Name() + "'s return value differs", false
		}
	}

	return "", true
}

func (fi *FunctionInterface) isImplementedByBuiltinFunction(fn prototypes.BuiltinFunction) (string, bool) {
	// generate dummy args
	dummyArgs := make([]values.Value, len(fi.args))

	for i, arg := range fi.args {
		if arg.constraint == nil {
			dummyArgs[i] = values.NewVoid(arg.Context())
		} else if globalProto := GetBuiltinPrototype(arg.constraint.Name()); globalProto != nil {

			dummyArgs[i] = values.NewNull(globalProto, arg.Context())
		} else {
			return fi.Name() + " differs", false
		}
	}

	res := fn.CheckArgs(dummyArgs)
	if !res {
		return fi.Name() + " differs", res
	} else {
		return "", true
	}

	// TODO: take fi.ret into account
	// XXX: this is actually quite difficult right now
}

// for builtin interfaces
func PrototypeImplements(proto_ values.Prototype, interf_ values.Interface) (string, bool) {
	proto, ok := proto_.(*Class)
	if !ok {
		return "not a class", false
	}

	interf, ok := interf_.(*prototypes.BuiltinInterface)
	if !ok {
		panic("not a builtin interface")
	}

	for name, member := range interf.Members {
		clMember := proto.FindMember(name, false, prototypes.IsSetter(member))
		if clMember == nil {
			return name + " missing", false
		}

		if clMember.Role() != member.Role() {
			return name + " roles differ", false
		}

		// we dont have access to the stack so we need another way
		fi := clMember.function.Interface()

		args := make([]interface{}, 0)
		for _, arg := range fi.args {
			args = append(args, arg.ConstraintName())
		}

		if !member.Check(args) {
			return name + " arg types differ", false
		}

		retName := ""
		if fi.ret != nil {
			retName = fi.ret.Name()
		}
		if !member.CheckRetType(retName) {
			return name + " return type differs", false
		}
	}

	// TODO: take fi.ret into account

	return "", true
}

var _PrototypeImplementsOk = prototypes.RegisterPrototypeImplements(PrototypeImplements)
