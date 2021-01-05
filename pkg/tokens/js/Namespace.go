package js

import (
	"github.com/computeportal/wtsuite/pkg/tokens/patterns"
)

type Namespace interface {
	NewBlockNamespace() Namespace
	NewFunctionNamespace() Namespace

	CurrentFunctionNamespace() Namespace

	UniversalName(v Variable, name string) error
	ClassName(v Variable) error // classnames must be the same project wide!
	FunctionName(v Variable)
	ArgName(v Variable)
	LetName(v Variable)
	VarName(v Variable)

	LibName(v Variable, name string) error

	HasName(new string) bool
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

func (ns *NamespaceData) UniversalName(v Variable, name string) error {
	return ns.OriginalName(v, name)
}

func (ns *NamespaceData) ClassName(v Variable) error {
	ns.VarName(v)
	return nil
}

func (ns *NamespaceData) OriginalName(v Variable, name string) error {
	if ns.HasVar(v) {
		// assumed to already have been succesfull before
		return nil
	}

	if ns.HasName(name) {
		otherVar := ns.nameVars[name]
		errCtx := v.Context()

		err := errCtx.NewError("Error: name '" + name + "' must be unique project wide!")
		otherCtx := otherVar.Context()
		err.AppendContextString("Info: previous usage of name", otherCtx)
		panic(err)
		return err
	}

	ns.varNames[v] = name
	ns.nameVars[name] = v
	v.Rename(name)

	return nil
}

func (ns *NamespaceData) FunctionName(v Variable) {
	ns.VarName(v)
}

func (ns *NamespaceData) ArgName(v Variable) {
	// update: can also be called in catch
	/*if !ns.isFunction {
		panic("can only be called immediately in function")
	}*/

	ns.VarName(v)
}

func (ns *NamespaceData) LetName(v Variable) {
	if ns.HasVar(v) {
		// already handled before, eg. by export
		return
	}

	ng := patterns.NewNameGenerator(true, v.Name())

	for true {
		new := ng.GenName()

		if !ns.HasName(new) {
			ns.nameVars[new] = v
			ns.varNames[v] = new
			v.Rename(new)

			return
		}
	}

	panic("impossible")
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
		new := ng.GenName()

		if !fns.HasName(new) && !ns.HasName(new) {
			fns.varNames[v] = new
			fns.nameVars[new] = v
			v.Rename(new)

			return
		}
	}

	panic("impossible")
}

func (ns *NamespaceData) LibName(v Variable, name string) error {
	if name == "" {
		panic("cant be empty")
	}

	return ns.OriginalName(v, name)
}

func (ns *NamespaceData) HasName(new string) bool {
	if _, ok := ns.nameVars[new]; ok {
		return true
	}

	if ns.parent != nil {
		return ns.parent.HasName(new)
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
