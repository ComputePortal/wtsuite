package macros

import (
  "fmt"
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/js"

	"github.com/computeportal/wtsuite/pkg/tokens/js/prototypes"
	"github.com/computeportal/wtsuite/pkg/tokens/js/values"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type BigIntCall struct {
	Macro
}

func NewBigIntCall(args []js.Expression, ctx context.Context) (js.Expression, error) {
  if len(args) != 1 {
    errCtx := ctx
    return nil, errCtx.NewError(fmt.Sprintf("Error: expected 1 argument, got %d", len(args)))
  }

	return &BigIntCall{newMacro(args, ctx)}, nil
}

func (m *BigIntCall) Dump(indent string) string {
	return indent + "BigIntCall(...)"
}

func (m *BigIntCall) WriteExpression() string {
	// XXX: should everything be wrapped in additional parentheses?
	var b strings.Builder

	b.WriteString("BigInt(")
	b.WriteString(m.args[0].WriteExpression())
	b.WriteString(")")

	return b.String()
}

func (m *BigIntCall) EvalExpression() (values.Value, error) {
	ctx := m.Context()
	args, err := m.evalArgs()
	if err != nil {
		return nil, err
	}

  if !(prototypes.IsString(args[0]) || prototypes.IsInt(args[0])) {
    errCtx := args[0].Context()
    return nil, errCtx.NewError("Error: expected String or Int as argument")
  }

	return prototypes.NewBigInt(ctx), nil
}
