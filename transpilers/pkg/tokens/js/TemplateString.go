package js

import (
	"strings"

	"./prototypes"
	"./values"

	"../context"
)

type TemplateString struct {
	lit []*LiteralString
	ins []Expression
	LiteralData
}

func NewTemplateString(exprs []Expression, ctx context.Context) (*TemplateString, error) {
	lit := make([]*LiteralString, 0)
	ins := make([]Expression, 0)

	for i, expr := range exprs {
		if i%2 == 0 {
			l, ok := expr.(*LiteralString)
			if !ok {
				errCtx := expr.Context()
				return nil, errCtx.NewError("Error: expected literal String")
			}
			lit = append(lit, l)
		} else {
			ins = append(ins, expr)
		}
	}

	if len(lit) != len(ins) && len(lit) != len(ins)+1 {
		return nil, ctx.NewError("Error expected at least as many literal parts as inserted parts")
	}

	return &TemplateString{lit, ins, newLiteralData(ctx)}, nil
}

func (t *TemplateString) Dump(indent string) string {
	var b strings.Builder

	b.WriteString(indent)
	b.WriteString("TemplateString\n")

	for i, l := range t.lit {
		b.WriteString(l.Dump(indent + "  "))

		if i < len(t.ins) {
			b.WriteString(t.ins[i].Dump(indent + "  "))
		}
	}

	return b.String()
}

func (t *TemplateString) WriteExpression() string {
	var b strings.Builder

	b.WriteString("`")

	for i, l := range t.lit {
		b.WriteString(l.Value())

		if i < len(t.ins) {
			b.WriteString("${")
			b.WriteString(t.ins[i].WriteExpression())
			b.WriteString("}")
		}
	}

	b.WriteString("`")

	return b.String()
}

func (t *TemplateString) ResolveExpressionNames(scope Scope) error {
	for _, ins := range t.ins {
		if err := ins.ResolveExpressionNames(scope); err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateString) EvalExpression(stack values.Stack) (values.Value, error) {
	for _, ins := range t.ins {
		v, err := ins.EvalExpression(stack)
		if err != nil {
			return nil, err
		}

		if !v.IsInstanceOf(prototypes.String, prototypes.Number, prototypes.Boolean) {
			errCtx := ins.Context()
			return nil, errCtx.NewError("Error: expected String(able)")
		}
	}

	if len(t.ins) == 0 && len(t.lit) == 1 {
		return prototypes.NewLiteralString(t.lit[0].value, t.Context()), nil
	} else {
		return prototypes.NewString(t.Context()), nil
	}
}

func (t *TemplateString) ResolveExpressionActivity(usage Usage) error {
	for _, ins := range t.ins {
		if err := ins.ResolveExpressionActivity(usage); err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateString) UniversalExpressionNames(ns Namespace) error {
	for _, ins := range t.ins {
		if err := ins.UniversalExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateString) UniqueExpressionNames(ns Namespace) error {
	for _, ins := range t.ins {
		if err := ins.UniqueExpressionNames(ns); err != nil {
			return err
		}
	}

	return nil
}

func (t *TemplateString) Walk(fn WalkFunc) error {
	for _, ins := range t.ins {
		if err := ins.Walk(fn); err != nil {
			return err
		}
	}

	return fn(t)
}
