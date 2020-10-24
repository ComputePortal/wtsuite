package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type BlobFromInstance struct {
	Macro
}

func NewBlobFromInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &BlobFromInstance{newMacro(args, ctx)}, nil
}

func (m *BlobFromInstance) Dump(indent string) string {
	return indent + "BlobFromInstance(...)"
}

func (m *BlobFromInstance) WriteExpression() string {
	var b strings.Builder

	b.WriteString(blobFromInstanceHeader.Name())
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *BlobFromInstance) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 1 {
		return nil, ctx.NewError("Error: expected 1 argument")
	}

	return prototypes.NewInstance(prototypes.Blob, ctx), nil
}

func (m *BlobFromInstance) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(blobFromInstanceHeader, m.Context())

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *BlobFromInstance) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(blobFromInstanceHeader, ns); err != nil {
		return err
	}

	return m.Macro.UniqueExpressionNames(ns)
}
