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
	Macro
}

func NewObjectToInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &ObjectToInstance{newToInstance(), newMacro(args, ctx)}, nil
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

	return b.String()
}

func (m *ObjectToInstance) EvalExpression() (values.Value, error) {
	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if len(args) != 2 {
		errCtx := m.Context()
		err := errCtx.NewError("Error: expected 2 arguments")
		return nil, err
	}

  objectCheck := prototypes.NewObject(nil, args[0].Context())
  if err := objectCheck.Check(args[0], args[0].Context()); err != nil {
    return nil, err
  }

	return args[1].EvalConstructor(nil, m.Context())
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
