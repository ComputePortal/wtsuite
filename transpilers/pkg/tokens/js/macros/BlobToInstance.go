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
	Macro
}

func NewBlobToInstance(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &BlobToInstance{newToInstance(), newMacro(args, ctx)}, nil
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

func (m *BlobToInstance) EvalExpression(stack values.Stack) (values.Value, error) {
	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if len(args) != 2 {
		errCtx := m.Context()
		err := errCtx.NewError("Error: expected 2 arguments")
		return nil, err
	}

	if !args[0].IsInstanceOf(prototypes.Blob) {
		errCtx := m.args[0].Context()
		return nil, errCtx.NewError("Error: expected Blob, got " + args[0].TypeName())
	}

	res, err := m.evalInstancePrototype(stack, args[1], m.Context())
	if err != nil {
		return nil, err
	}

	return prototypes.NewResolvedPromise(res, m.Context())
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