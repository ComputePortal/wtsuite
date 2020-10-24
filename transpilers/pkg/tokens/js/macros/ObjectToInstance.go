package macros

import (
	"strings"

	"../"

	"../values"

	"../../context"
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

func (m *ObjectToInstance) EvalExpression(stack values.Stack) (values.Value, error) {

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 2 {
		errCtx := m.Context()
		err := errCtx.NewError("Error: expected 2 arguments")
		return nil, err
	}

	// XXX: arg[0] can be any, or should it at least be Object type?

	return m.evalInstancePrototype(stack, args[1], m.Context())
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
