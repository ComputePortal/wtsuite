package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type PostMacro struct {
	ToInstance
	Macro
}

func newPostMacro(args []js.Expression, ctx context.Context) PostMacro {
	return PostMacro{newToInstance(), newMacro(args, ctx)}
}

func (m *PostMacro) writeExpression(fnName string) string {
	var b strings.Builder

	b.WriteString(fnName)
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(",")
	b.WriteString(m.args[1].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *PostMacro) evalExpression(msg values.Value,
	classValue values.Value) (values.Value, error) {
	if !isAnObject(msg) {
		errCtx := m.Context()
		return nil,
			errCtx.NewError("Error: expected Object or instance of class that extends Object for argument 2, got " +
				msg.TypeName())
	}

	resolveValue, err := classValue.EvalConstructor(nil, classValue.Context())
	if err != nil {
		context.AppendContextString(err, "Info: needed here", m.Context())
		return nil, err
	}

  proto := values.GetPrototype(resolveValue)
  if proto == nil {
    panic("expected instance of class")
  }

  if !proto.IsUniversal() {
    errCtx := m.Context()
    return nil, errCtx.NewError("Error: class " + proto.Name() + " is not universal (hint: use 'universe'")
  }

  return prototypes.NewPromise(resolveValue, m.Context()), nil
}
