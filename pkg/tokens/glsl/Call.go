package glsl

import (
	"strings"

	"github.com/computeportal/wtsuite/pkg/tokens/context"
)

type Call struct {
	lhs  Expression
	args []Expression
	TokenData
}

func NewCall(lhs Expression, args []Expression, ctx context.Context) *Call {
	return &Call{lhs, args, TokenData{ctx}}
}

// returns empty string if lhs is not *VarExpression
func (t *Call) Name() string {
	if ve, ok := t.lhs.(*VarExpression); ok {
		return ve.Name()
	} else {
		return ""
	}
}

func (t *Call) Args() []Expression {
	return t.args
}

func (t *Call) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("Call\n")

	b.WriteString(t.lhs.Dump(indent + "  "))

	for _, arg := range t.args {
		b.WriteString(arg.Dump(indent + "( "))
	}

	return b.String()
}

func (t *Call) WriteExpression() string {
	var b strings.Builder

	b.WriteString(t.lhs.WriteExpression())

	b.WriteString("(")

	for i, arg := range t.args {
		b.WriteString(arg.WriteExpression())

		if i < len(t.args)-1 {
			b.WriteString(",")
		}
	}

	b.WriteString(")")

	return b.String()
}

func (t *Call) WriteStatement(usage Usage, indent string, nl string, tab string) string {
	return indent + t.WriteExpression()
}
