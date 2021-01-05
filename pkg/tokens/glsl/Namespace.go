package glsl

import (
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type Namespace interface {
	NewBlockNamespace() Namespace
	NewFunctionNamespace() Namespace

	CurrentFunctionNamespace() Namespace

	FunctionName(v Variable)
	ArgName(v Variable)
	VarName(v Variable)
  OrigName(v Variable) error

	HasName(newName string) bool
	HasVar(v Variable) bool
}

type NamespaceData struct {
	parent Namespace

	isFunction bool

	varNames map[Variable]string // variable -> new (new is also stored in variable itself)
	nameVars map[string]Variable // new -> variable
}

func newNamespace(parent Namespace, isFunction bool) Namespace {
	return &NamespaceData{parent, isFunction, make(map[Variable]string), make(map[string]Variable)}
}

func NewNamespace(parent Namespace, isFunction bool) Namespace {
	return newNamespace(parent, isFunction)
}

func (ns *NamespaceData) NewBlockNamespace() Namespace {
	return newNamespace(ns, false)
}

func (ns *NamespaceData) NewFunctionNamespace() Namespace {
	return newNamespace(ns, true)
}

func (ns *NamespaceData) CurrentFunctionNamespace() Namespace {
	if ns.isFunction || ns.parent == nil {
		return ns
	} else {
		return ns.parent.CurrentFunctionNamespace()
	}
}

func (ns *NamespaceData) OrigName(v Variable) error {
  name := v.Name()

  if ns.HasVar(v) {
    // assumed already ok before
    return nil
  }

	if ns.HasName(name) {
		otherVar := ns.nameVars[name]
		errCtx := v.Context()

		err := errCtx.NewError("Error: name '" + name + "' must be unique")
		otherCtx := otherVar.Context()
		err.AppendContextString("Info: previous usage of name", otherCtx)
		return err
	}

	ns.varNames[v] = name
	ns.nameVars[name] = v

	return nil
}

func (ns *NamespaceData) FunctionName(v Variable) {
	ns.VarName(v)
}

func (ns *NamespaceData) ArgName(v Variable) {
	ns.VarName(v)
}

func (ns *NamespaceData) VarName(v Variable) {
	if ns.HasVar(v) {
		// already handled before, eg. by export
		return
	}

	ng := patterns.NewNameGenerator(true, v.Name())

	fns_ := ns.CurrentFunctionNamespace()

	fns, ok := fns_.(*NamespaceData)
	if !ok {
		panic("unexpected")
	}

	for true {
		newName := ng.GenName()

		if !fns.HasName(newName) && !ns.HasName(newName) {
			fns.varNames[v] = newName
			fns.nameVars[newName] = v
			v.Rename(newName)

			return
		}
	}

	panic("impossible")
}

func (ns *NamespaceData) HasName(newName string) bool {
	if _, ok := ns.nameVars[newName]; ok {
		return true
	}

	if ns.parent != nil {
		return ns.parent.HasName(newName)
	}

	return false
}

func (ns *NamespaceData) HasVar(v Variable) bool {
	if name, ok := ns.varNames[v]; ok {
		if name != v.Name() {
			panic("something went wrong")
		}

		return true
	}

	if ns.parent != nil {
		return ns.parent.HasVar(v)
	}

	return false
}
