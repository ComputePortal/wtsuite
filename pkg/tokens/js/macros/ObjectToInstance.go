package macros

import (
	"strings"

  "../prototypes"

	"../values"

	"../../context"
	"../../js"
)

type ObjectToInstance struct {
	ToInstance
}

func NewObjectToInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

  interfExpr, err := getTypeExpression(args[1])
  if err != nil {
    return nil, err
  }

  return &ObjectToInstance{newToInstance(args[0:1], interfExpr, ctx)}, nil
}

func (m *ObjectToInstance) Dump(indent string) string {
	return indent + "ObjectToInstance(...)"
}

func (m *ObjectToInstance) WriteExpression() string {
	var b strings.Builder

	b.WriteString(objectToInstanceHeader.Name())
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return m.wrapWithCheckType(b.String())
}

func (m *ObjectToInstance) EvalExpression() (values.Value, error) {
  args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

  if !prototypes.IsObject(args[0]) {
    errCtx := args[0].Context()
    return nil, errCtx.NewError("Error: expected Object, got " + args[0].TypeName())
  }

	return values.NewInstance(m.interf, m.Context()), nil
}

func (m *ObjectToInstance) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(objectToInstanceHeader, m.Context())

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *ObjectToInstance) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(objectToInstanceHeader, ns); err != nil {
		return err
	}

	return m.Macro.UniqueExpressionNames(ns)
}
