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
}

func newPostMacro(args []js.Expression, interfExpr *js.TypeExpression, ctx context.Context) PostMacro {
	return PostMacro{newToInstance(args, interfExpr, ctx)}
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

func (m *PostMacro) evalExpression(msg values.Value) (values.Value, error) {
	if !isAnObject(msg) {
		errCtx := m.Context()
		return nil, errCtx.NewError("Error: expected Object or instance of class that extends Object for argument 2, got " + msg.TypeName())
	}

  res := values.NewInstance(m.interf, m.Context())

  return prototypes.NewPromise(res, m.Context()), nil
}
