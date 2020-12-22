package directives

import (
	"fmt"
	"strconv"
	"strings"

	"../files"

	tokens "../tokens/html"

	"../tokens/js"
	"../tokens/js/prototypes"
	"../tokens/js/values"

	"../tokens/context"
	"../tree"
	"../tree/scripts"
)

type ViewVariable interface {
	js.Variable
	Write() string
}

// use this to import ui variables into script files
type ViewModule struct {
	scope    *ScopeData
	imported map[string]ViewVariable
	path     string
	ctx      context.Context
}

type ViewFileScript struct {
	module *ViewModule
}

type ViewVariableData struct {
	name string
	ctx  context.Context
}

func newViewVariableData(name string, ctx context.Context) ViewVariableData {
	return ViewVariableData{name, ctx}
}

// used as fallback when js.TARGET != "browser"
type ViewNull struct {
	ViewVariableData
}

type ViewString struct {
	value string
	ViewVariableData
}

type ViewInt struct {
	value int
	ViewVariableData
}

type ViewBool struct {
	value bool
	ViewVariableData
}

type ViewNumber struct {
	value float64
	ViewVariableData
}

type ViewColor struct {
	r, g, b, a float64 // between 0.0 and 1.0, returns a small literal array in js
	ViewVariableData
}

type ViewArray struct {
	values []ViewVariable
	ViewVariableData
}

type ViewObject struct {
	values map[string]ViewVariable
	ViewVariableData
}

func NewViewFileScript(absPath string, caller string) (scripts.FileScript, error) {
	if js.TARGET == "browser" || js.TARGET == "all" {
		fileScope, rootNode, err := BuildFile(absPath, caller, false)
		if err != nil {
			return nil, err
		}

		return &ViewFileScript{&ViewModule{fileScope, make(map[string]ViewVariable), absPath, rootNode.Context()}}, nil
	} else {
		// returns allnull on every request
		return &ViewFileScript{&ViewModule{nil, make(map[string]ViewVariable), absPath, context.NewDummyContext()}}, nil
	}
}

func (m *ViewModule) isDummy() bool {
	return m.scope == nil
}

func (m *ViewModule) tokenToViewVariable(t tokens.Token, name string, ctx context.Context) (ViewVariable, error) {
	// for nested value the name is not actually used
	switch v := t.(type) {
	case *tokens.String:
		return &ViewString{v.Value(), newViewVariableData(name, ctx)}, nil
	case *tokens.Bool:
		return &ViewBool{v.Value(), newViewVariableData(name, ctx)}, nil
	case *tokens.Int:
		return &ViewInt{v.Value(), newViewVariableData(name, ctx)}, nil
	case *tokens.Float:
		if v.Unit() != "" {
			return nil, ctx.NewError("Error: can't import united float from view")
		}
		return &ViewNumber{v.Value(), newViewVariableData(name, ctx)}, nil
	case *tokens.Color:
		r, g, b, a := v.FloatValues()
		return &ViewColor{r, g, b, a, newViewVariableData(name, ctx)}, nil
	case *tokens.List:
		vts := v.GetTokens()
		arr := make([]ViewVariable, len(vts))
		for i, nested := range vts {
			var err error
			arr[i], err = m.tokenToViewVariable(nested, name, ctx)
			if err != nil {
				return nil, err
			}
		}

		return &ViewArray{arr, newViewVariableData(name, ctx)}, nil
	case *tokens.StringDict:
		obj := make(map[string]ViewVariable)
		if err := v.Loop(func(key *tokens.String, nested_ tokens.Token, last bool) error {
			nested, err := m.tokenToViewVariable(nested_, name, ctx)
			if err != nil {
				return err
			}

			obj[key.Value()] = nested
			return nil
		}); err != nil {
			return nil, err
		}

		return &ViewObject{obj, newViewVariableData(name, ctx)}, nil
	default:
		return nil, ctx.NewError("Error: to importable from view")
	}
}

func (m *ViewModule) GetExportedVariable(gs js.GlobalScope, name string,
	ctx context.Context) (js.Variable, error) {
	if prev, ok := m.imported[name]; ok {
		return prev, nil
	} else {
		if m.isDummy() {
			v := &ViewNull{newViewVariableData(name, ctx)}
			m.imported[name] = v
			return v, nil
		} else {
			if !m.scope.HasVar(name) {
				return nil, ctx.NewError("Error: " + name + " not found in " + files.Abbreviate(m.path))
			}

			v_ := m.scope.GetVar(name)
			if !v_.Exported {
				return nil, ctx.NewError("Error: " + name + " not exported by " + files.Abbreviate(m.path))
			}

			v, err := m.tokenToViewVariable(v_.Value, name, ctx)
			if err != nil {
				return nil, err
			}

			m.imported[name] = v

			return v, nil
		}
	}
}

func (m *ViewModule) Context() context.Context {
	return m.ctx
}

func (s *ViewFileScript) Write() (string, error) {
	if s.module.isDummy() {
		return "", nil
	}

	var b strings.Builder

	for _, v := range s.module.imported {
		b.WriteString("var ")
		b.WriteString(v.Name())
		b.WriteString("=")
		b.WriteString(v.Write())
		b.WriteString(";")
		b.WriteString(tree.NL)
	}

	return b.String(), nil
}

func (s *ViewFileScript) Dependencies() []string {
	// cache dependencies should've been created during importExport
	return []string{}
}

// no need to add to the scope
func (s *ViewFileScript) ResolveNames(scope js.GlobalScope) error {
	return nil
}

func (s *ViewFileScript) EvalTypes() error {
	return nil
}

func (s *ViewFileScript) ResolveActivity(usage js.Usage) error {
	return nil
}

func (s *ViewFileScript) UniqueEntryPointNames(ns js.Namespace) error {
	return nil
}

func (s *ViewFileScript) UniversalNames(ns js.Namespace) error {
	return nil
}

func (s *ViewFileScript) UniqueNames(ns js.Namespace) error {
	for _, v := range s.module.imported {
		ns.VarName(v)
	}

	return nil
}

func (s *ViewFileScript) Module() js.Module {
	return s.module
}

func (s *ViewFileScript) Path() string {
	return s.module.path
}

func (s *ViewFileScript) Walk(fn func(p string, obj interface{}) error) error {
  // TODO: should we walk the ViewModule? Or is this part of htmlpp refactoring?
  return nil
}

func (v *ViewVariableData) Context() context.Context {
	return v.ctx
}

func (v *ViewVariableData) Name() string {
	return v.name
}

func (v *ViewVariableData) Constant() bool {
	return true
}

func (v *ViewVariableData) SetConstant() {
}

func (v *ViewVariableData) Rename(newName string) {
	v.name = newName
}

func (v *ViewVariableData) GetObject() interface{} {
	return nil
}

func (v *ViewVariableData) SetObject(interface{}) {
}

func (v *ViewVariableData) SetValue(_ values.Value) {
  panic("cant be set")
}

func (v *ViewNull) Dump(indent string) string {
	return indent + "ViewNull"
}

func (v *ViewNull) GetValue() values.Value {
	return values.NewAll(v.Context())
}

func (v *ViewNull) Write() string {
	return "null"
}

func (v *ViewString) Dump(indent string) string {
	return indent + "ViewString(" + v.value + ")"
}

func (v *ViewString) GetValue() values.Value {
	return prototypes.NewLiteralString(v.value, v.Context())
}

func (v *ViewString) Write() string {
	return "\"" + v.value + "\""
}

func (v *ViewInt) Dump(indent string) string {
	return "ViewInt(" + strconv.Itoa(v.value) + ")"
}

func (v *ViewInt) GetValue() values.Value {
	return prototypes.NewLiteralInt(v.value, v.Context())
}

func (v *ViewInt) Write() string {
	return strconv.Itoa(v.value)
}

func (v *ViewBool) Dump(indent string) string {
	if v.value {
		return "ViewBool(true)"
	} else {
		return "ViewBool(false)"
	}
}

func (v *ViewBool) GetValue() values.Value {
	return prototypes.NewLiteralBoolean(v.value, v.Context())
}

func (v *ViewBool) Write() string {
	if v.value {
		return "true"
	} else {
		return "false"
	}
}

func (v *ViewNumber) Dump(indent string) string {
	return fmt.Sprintf("ViewNumber(%g)", v.value)
}

func (v *ViewNumber) GetValue() values.Value {
	return prototypes.NewNumber(v.Context())
}

func (v *ViewNumber) Write() string {
	return fmt.Sprintf("%g", v.value)
}

func (v *ViewColor) Dump(indent string) string {
	return fmt.Sprintf("ViewColor(%g,%g,%g,%g)", v.r, v.g, v.b, v.a)
}

func (v *ViewColor) GetValue() values.Value {
	ctx := v.Context()

	return prototypes.NewArray(prototypes.NewNumber(ctx), ctx)
}

func (v *ViewColor) Write() string {
	return fmt.Sprintf("[%g,%g,%g,%g]", v.r, v.g, v.b, v.a)
}

func (v *ViewArray) Dump(indent string) string {
	var b strings.Builder
	b.WriteString(indent)
	b.WriteString("ViewArray\n")
	for _, item := range v.values {
		b.WriteString(item.Dump(indent + "  "))
		b.WriteString("\n")
	}

	return b.String()
}

func (v *ViewArray) GetValue() values.Value {
	items := make([]values.Value, len(v.values))

	for i, item := range v.values {
		items[i] = item.GetValue()
	}

	return prototypes.NewArray(values.CommonValue(items, v.Context()), v.Context())
}

func (v *ViewArray) Write() string {
	var b strings.Builder

	b.WriteString("[")

	for i, item := range v.values {
		b.WriteString(item.Write())
		if i < len(v.values)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString("]")

	return b.String()
}

func (v *ViewObject) Dump(indent string) string {
	var b strings.Builder
	b.WriteString(indent + "ViewObject\n")

	for k, item := range v.values {
		b.WriteString(indent)
		b.WriteString(k)
		b.WriteString(":\n")

		b.WriteString(item.Dump(indent + "  "))
		b.WriteString("\n")
	}

	return b.String()
}

func (v *ViewObject) GetValue() values.Value {
	items := make(map[string]values.Value)

	for k, item := range v.values {
		items[k] = item.GetValue()
	}

	return prototypes.NewObject(items, v.Context())
}

func (v *ViewObject) Write() string {
	var b strings.Builder

	b.WriteString("{")

	i := 0
	for k, item := range v.values {
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString(item.Write())

		if i < len(v.values)-1 {
			b.WriteString(",")
		}

		i++
	}
	b.WriteString("}")

	return b.String()
}

var _newViewFileScriptOk = scripts.SetNewViewFileScript(NewViewFileScript)

// in jspp the directives module isn't automatically included, so we need to force (which includes the above lines actually, so this happens twice)
func ForceNewViewFileScriptRegistration() {
	scripts.SetNewViewFileScript(NewViewFileScript)
}
