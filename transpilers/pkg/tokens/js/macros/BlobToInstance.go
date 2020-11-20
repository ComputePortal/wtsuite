package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type BlobToInstance struct {
	ToInstance
}

func NewBlobToInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	if len(args) != 2 {
		return nil, ctx.NewError("Error: expected 2 arguments")
	}

  interfExpr, err := getTypeExpression(args[1])
  if err != nil {
    return nil, err
  }

  return &BlobToInstance{newToInstance(args[0:1], interfExpr, ctx)}, nil
}

func (m *BlobToInstance) Dump(indent string) string {
	return indent + "BlobToInstance(...)"
}

func (m *BlobToInstance) WriteExpression() string {
	var b strings.Builder

	b.WriteString(blobToInstanceHeader.Name())
	b.WriteString("(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *BlobToInstance) EvalExpression() (values.Value, error) {
  args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

	if !prototypes.IsBlob(args[0]) {
		errCtx := m.args[0].Context()
		return nil, errCtx.NewError("Error: expected Blob, got " + args[0].TypeName())
	}

  res := values.NewInstance(m.interf, m.Context())

	return prototypes.NewPromise(res, m.Context()), nil
}

func (m *BlobToInstance) ResolveExpressionActivity(usage js.Usage) error {
	ResolveHeaderActivity(blobToInstanceHeader, m.Context())

	return m.Macro.ResolveExpressionActivity(usage)
}

func (m *BlobToInstance) UniqueExpressionNames(ns js.Namespace) error {
	if err := UniqueHeaderNames(blobToInstanceHeader, ns); err != nil {
		return err
	}

	return m.Macro.UniqueExpressionNames(ns)
}
