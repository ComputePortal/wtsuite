package macros

import (
	"strings"

	"../"

	"../prototypes"
	"../values"

	"../../context"
)

type DocumentGetElementById struct {
	ToInstance
	Macro
}

func NewDocumentGetElementById(args []js.Expression, ctx context.Context) (js.Expression, error) {
	return &DocumentGetElementById{newToInstance(), newMacro(args, ctx)}, nil
}

func (m *DocumentGetElementById) Dump(indent string) string {
	return indent + "DocumentGetElementById(...)"
}

func (m *DocumentGetElementById) WriteExpression() string {
	var b_ strings.Builder

	b_.WriteString("document.getElementById(")
	b_.WriteString(m.args[0].WriteExpression())
	b_.WriteString(")")

	if COMPACT {
		return b_.String()
	} else {
		var b strings.Builder

		b.WriteString("(()=>{let e=")
		b.WriteString(b_.String())
		b.WriteString(";if(e===null){throw new Error('element #'+")
		b.WriteString(m.args[0].WriteExpression())
		b.WriteString("+' not found');};")

		if len(m.args) == 2 {
			b.WriteString("if(!(e instanceof ")
			b.WriteString(m.args[1].WriteExpression())
			b.WriteString(")){throw new Error('element #'+")
			b.WriteString(m.args[0].WriteExpression())
			b.WriteString("+' is not a ")
			b.WriteString(m.args[1].WriteExpression())
			b.WriteString(" (but a: '+e.constructor.name+')');};")
		}
		b.WriteString("return e})()")

		return b.String()
	}
}

func (m *DocumentGetElementById) EvalExpression(stack values.Stack) (values.Value, error) {
	ctx := m.Context()

	args, err := m.evalArgs(stack)
	if err != nil {
		return nil, err
	}

	if !args[0].IsInstanceOf(prototypes.String) {
		return nil, ctx.NewError("Error: expected String for first argument, got " + args[0].TypeName())
	}

	if len(args) == 1 {
		return prototypes.NewInstance(prototypes.HTMLElement, ctx), nil
	} else if len(args) == 2 {
		outClass, ok := args[1].GetClassPrototype()
		if !ok {
			return nil, ctx.NewError("Error: second argument is not a class")
		}

		if !outClass.HasAncestor(prototypes.HTMLElement) {
			return nil, ctx.NewError("Error: " + outClass.Name() + " doesn't inherit from HTMLElement")
		}

		return prototypes.NewInstance(outClass, ctx), nil
	} else {
		return nil, ctx.NewError("Error: expected 1 or 2 arguments")
	}
}
