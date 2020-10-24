package macros

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type SyntaxTreeInfo struct {
	evalCount int
	Macro
}

func NewSyntaxTreeInfo(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &SyntaxTreeInfo{0, newMacro(args, ctx)}, nil
}

func (m *SyntaxTreeInfo) Dump(indent string) string {
	return indent + "SyntaxTreeInfo(...)"
}

func (m *SyntaxTreeInfo) WriteExpression() string {
	return m.args[0].WriteExpression()
}

func (m *SyntaxTreeInfo) dumpValue(v values.Value) string {
	var b strings.Builder
	b.WriteString(v.TypeName())
	switch {
	case v.IsInstanceOf(prototypes.Int):
		if intVal, ok := v.LiteralIntValue(); ok {
			b.WriteString(fmt.Sprintf("(%d)", intVal))
		}
	case v.IsInstanceOf(prototypes.Boolean):
		if boolVal, ok := v.LiteralBooleanValue(); ok {
			if boolVal {
				b.WriteString("(true)")
			} else {
				b.WriteString("(false)")
			}
		}
	case v.IsInstanceOf(prototypes.String):
		if strVal, ok := v.LiteralStringValue(); ok {
			b.WriteString("(")
			b.WriteString(strVal)
			b.WriteString(")")
		}
	case v.IsNull():
		b.WriteString("(null)")
	}

	return b.String()
}

func (m *SyntaxTreeInfo) contextMessage() string {
	if len(m.args) == 2 {
		if lit, ok := m.args[1].(*js.LiteralString); ok {
			return lit.Value()
		} else {
			panic("should've been caught before")
		}
	} else {
		return ""
	}
}

func (m *SyntaxTreeInfo) ResolveExpressionNames(scope js.Scope) error {
	if len(m.args) != 1 && len(m.args) != 2 {
		errCtx := m.Context()
		return errCtx.NewError("Error: expected 1 or 2 arguments")
	}

	if len(m.args) == 2 {
		if _, ok := m.args[1].(*js.LiteralString); !ok {
			errCtx := m.args[1].Context()
			return errCtx.NewError("Error: expected literal string as second argument")
		}
	}

	var v js.Variable = nil
	switch expr := m.args[0].(type) {
	case *js.VarExpression:
		v = expr.GetVariable()
	}

	prefix := m.contextMessage()
	if prefix == "" {
		prefix = "SYNTAX_TREE_VAR_INFO"
	}

	if v != nil {
		fmt.Fprintf(os.Stdout, "#%s: %p\n", prefix, v)
	} else {
		fmt.Fprintf(os.Stdout, "#%s: NA\n", prefix)
	}

	return m.Macro.ResolveExpressionNames(scope)
}

func (m *SyntaxTreeInfo) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	prefix := m.contextMessage()
	if prefix == "" {
		prefix = "SYNTAX_TREE_VALUE_INFO"
	}
	if len(args) == 2 {
		str, ok := args[1].LiteralStringValue()
		if !ok {
			errCtx := ctx
			return nil, errCtx.NewError("Error: expected literal string for argument 2")
		}
		prefix = str
	} else if len(args) != 1 {
		errCtx := ctx
		return nil, errCtx.NewError("Error: expected 1 or 2 argument")
	}

	arg := args[0]
	arg = values.UnpackContextValue(arg)

	suffix := ""
	if _, isMulti := arg.(*values.Multi); isMulti {
		suffix = "(multi)"
	} else {
		suffix = reflect.TypeOf(arg).String()
	}

	fmt.Fprintf(os.Stdout, "#%s: %s @ %p %s\n", prefix, m.dumpValue(arg), arg, suffix)

	if m.evalCount > 1999 {
		panic("too many, there is something wrong with the compilation")
	}

	m.evalCount += 1

	return args[0], nil
}
