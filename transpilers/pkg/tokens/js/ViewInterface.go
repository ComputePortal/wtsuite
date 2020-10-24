package js

import (
	"sort"
	"strings"

	"./values"

	"../context"
	"../raw"

	"../../files"
)

// to avoid cyclical imports
var _typeExpressionBuilder func(str string) *TypeExpression = nil

func SetTypeExpressionBuilder(fn func(string) *TypeExpression) bool {
	_typeExpressionBuilder = fn
	return true
}

type ViewInterfaceTypeValue struct {
	Type  string // eg. Int
	Value string // just for printing
}

type ViewInterfaceElement struct {
	Type   string
	HTML   string
	States map[string][]string
}

type ViewInterface struct {
	Path  string
	URL   string
	Vars  map[string]ViewInterfaceTypeValue
	Defs  map[string]ViewInterfaceTypeValue // value is the js constructor, type is the js nested type
	Elems map[string]ViewInterfaceElement   // rhs is the js nested type

	vars  map[string]*TypeExpression
	defs  map[string]*TypeExpression
	elems map[string]*TypeExpression
}

// path is absPath of ui file
func NewViewInterface(path, url string) *ViewInterface {
	return &ViewInterface{
		path,
		url,
		make(map[string]ViewInterfaceTypeValue),
		make(map[string]ViewInterfaceTypeValue),
		make(map[string]ViewInterfaceElement),
		make(map[string]*TypeExpression),
		make(map[string]*TypeExpression),
		make(map[string]*TypeExpression),
	}
}

func (vif *ViewInterface) IsSame(other *ViewInterface) bool {
	if len(vif.Vars) != len(other.Vars) {
		return false
	}

	for name, typeValue := range vif.Vars {
		if otherTypeValue, ok := other.Vars[name]; ok {
			if otherTypeValue.Type != typeValue.Type {
				return false
			}
		} else {
			return false
		}
	}

	if len(vif.Defs) != len(other.Defs) {
		return false
	}

	for name, typeValue := range vif.Defs {
		if otherTypeValue, ok := other.Defs[name]; ok {
			if otherTypeValue.Type != typeValue.Type {
				return false
			}
		} else {
			return false
		}
	}

	if len(vif.Elems) != len(other.Elems) {
		return false
	}

	for name, el := range vif.Elems {
		if otherEl, ok := other.Elems[name]; ok {
			if otherEl.Type != el.Type {
				return false
			} else if otherEl.HTML != el.HTML {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (vif *ViewInterface) AddVar(name string, t, v string) {
	vif.Vars[name] = ViewInterfaceTypeValue{t, v}
}

func (vif *ViewInterface) AddElem(id string, t string, html string, states map[string][]string) {
	vif.Elems[id] = ViewInterfaceElement{
		t,
		html,
		states,
	}
}

func (vif *ViewInterface) AddDef(name string, t string, v string) {
	vif.Defs[name] = ViewInterfaceTypeValue{t, v}
}

func (vif *ViewInterface) GetURL() string {
	return vif.URL
}

func (vif *ViewInterface) ResolveNames(scope Scope) error {
	// if constructed from gob, then vars/defs/elems will still be nil
	if vif.vars == nil {
		vif.vars = make(map[string]*TypeExpression)
	}
	if vif.defs == nil {
		vif.defs = make(map[string]*TypeExpression)
	}
	if vif.elems == nil {
		vif.elems = make(map[string]*TypeExpression)
	}

	for key, v := range vif.Vars {
		te := _typeExpressionBuilder(v.Type)

		if err := te.ResolveExpressionNames(scope); err != nil {
			return err
		}

		vif.vars[key] = te
	}

	for key, el := range vif.Elems {
		t := el.Type
		te := _typeExpressionBuilder(t)

		if err := te.ResolveExpressionNames(scope); err != nil {
			return err
		}

		vif.elems[key] = te
	}

	for key, v := range vif.Defs {
		te := _typeExpressionBuilder(v.Type)

		if err := te.ResolveExpressionNames(scope); err != nil {
			return err
		}

		vif.defs[key] = te
	}

	return nil
}

func (vif *ViewInterface) GetVarTypeInstance(stack values.Stack, name string,
	ctx context.Context) (values.Value, error) {
	if te, ok := vif.vars[name]; ok {
		return te.GenerateInstance(stack, ctx)
	} else {
		return nil, ctx.NewError("Error: '" + name + "' not a valid document var")
	}
}

func (vif *ViewInterface) writeValidElements() string {
	var b strings.Builder

	b.WriteString("\n")
	if len(vif.elems) == 0 {
		b.WriteString("\u001b[1mInfo: no elements available\u001b[0m\n")
	} else {
		b.WriteString("\u001b[1mInfo: valid elements\u001b[0m\n")

		// sort the elements alphabetically
		names := make([]string, 0)
		for name, _ := range vif.elems {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			b.WriteString("  ")
			b.WriteString(name)
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (vif *ViewInterface) GetElemStates(name string) map[string][]string {
	if el, ok := vif.Elems[name]; ok {
		return el.States
	} else {
		panic("element " + name + " not found?")
	}
}

func (vif *ViewInterface) GetElemTypeInstance(stack values.Stack, name string,
	ctx context.Context) (values.Value, error) {
	if te, ok := vif.elems[name]; ok {
		res, err := te.GenerateInstance(stack, ctx)
		if err != nil {
			return nil, err
		}

		res = values.UnpackContextValue(res)
		return res, nil
	} else {
		err := ctx.NewError("Error: " + files.Abbreviate(vif.Path) + "#" + name + " undefined")
		if VERBOSITY >= 1 {
			context.AppendString(err, vif.writeValidElements())
		}
		return nil, err
	}
}

func (vif *ViewInterface) GetHTML(stack values.Stack, name string,
	ctx context.Context) (string, error) {
	if el, ok := vif.Elems[name]; ok {
		return el.HTML, nil
	} else {
		err := ctx.NewError("Error: " + files.Abbreviate(vif.Path) + "#" + name + " undefined")
		if VERBOSITY >= 1 {
			context.AppendString(err, vif.writeValidElements())
		}
		return "", err
	}
}

func (vif *ViewInterface) GetDefTypeInstance(stack values.Stack, name string,
	ctx context.Context) (values.Value, error) {
	if te, ok := vif.defs[name]; ok {
		return te.GenerateInstance(stack, ctx)
	} else {
		return nil, ctx.NewError("Error: not a valid element constructor")
	}
}

func (vif *ViewInterface) WriteVariables() string {
	var b strings.Builder

	b.WriteString("{")

	for name, vt := range vif.Vars {
		// names can contain hyphens, so we need to surround them with quotes
		b.WriteString("\"")
		b.WriteString(name)
		b.WriteString("\":")
		b.WriteString(vt.Value)
		b.WriteString(",")
	}

	b.WriteString("}")

	return b.String()
}

func (vif *ViewInterface) WriteDefs() string {
	var b strings.Builder

	b.WriteString("{")

	for name, vt := range vif.Defs {
		// names can contain hyphens, so we need to surround the with quotes
		b.WriteString("\"")
		b.WriteString(name)
		b.WriteString("\":")
		b.WriteString(vt.Value)
		b.WriteString(",")
	}

	b.WriteString("}")

	return b.String()
}

func HashControl(fname string) string {
	return raw.ShortHash(fname)
}

func (vif *ViewInterface) IsElem(key string) bool {
	_, ok := vif.Elems[key]
	return ok
}
