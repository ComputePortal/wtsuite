package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type ObjectFromInstance struct {
	Macro
}

func NewObjectFromInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return &ObjectFromInstance{newMacro(args, ctx)}, nil
}

func (m *ObjectFromInstance) Dump(indent string) string {
	return indent + "ObjectFromInstance(...)"
}

func (m *ObjectFromInstance) WriteExpression() string {
	var b strings.Builder

	b.WriteString(objectFromInstanceHeader.Name())
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func isAnObject(v values.Value) bool {
  if values.IsAny(v) {
    return true
  } else if values.IsInstance(v) {
    if prototypes.IsNumber(v) || 
      prototypes.IsString(v) ||
      prototypes.IsArray(v) ||
      prototypes.IsBoolean(v) || 
      prototypes.IsTypedArray(v) ||
      prototypes.IsRegExp(v) {
      return false
    } else {
      return true
    }
	} else {
    return false
  }
}

func (m *ObjectFromInstance) EvalExpression() (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if !isAnObject(args[0]) {
		return nil,
			ctx.NewError("Error: expected Object or instance of class that extends Object for argument 1, got " +
				args[0].TypeName())
	}

  return prototypes.NewObject(nil, ctx), nil
}

func (m *ObjectFromInstance) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(objectFromInstanceHeader, m.Context())

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *ObjectFromInstance) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(objectFromInstanceHeader, ns); err != nil {
		return err
	}

	return m.Macro.UniqueExpressionNames(ns)
}
