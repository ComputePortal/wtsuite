package js

import (
	"./values"

	"../context"
)

type ForOf struct {
	await bool
	ForInOf
}

func NewForOf(await bool, varType VarType, lhs *VarExpression, rhs Expression,
	ctx context.Context) (*ForOf, error) {
	return &ForOf{await, newForInOf(varType, lhs, rhs, ctx)}, nil
}

func (t *ForOf) Dump(indent string) string {
	op := "of"
	if t.await {
		op += "await"
	}
	return t.ForInOf.dump(indent, op)
}

func (t *ForOf) WriteStatement(indent string) string {
	extra := ""
	if t.await {
		extra = "await"
	}
	return t.ForInOf.writeStatement(indent, extra, "of")
}

func (t *ForOf) EvalStatement(stack values.Stack) error {
	rhsValue, err := t.rhs.EvalExpression(stack)
	if err != nil {
		return err
	}

	valueCtx := t.lhs.Context()

	evalInner := func(v values.Value) error {
		subStack := NewBranchStack(stack)

		if t.varType == VAR {
			if err := stack.SetValue(t.lhs.GetVariable(), v, true,
				valueCtx); err != nil {
				return err
			}
		} else {
			if err := subStack.SetValue(t.lhs.GetVariable(), v, true,
				valueCtx); err != nil {
				return err
			}
		}

		return t.Block.EvalStatement(subStack)
	}

	if err := rhsValue.LoopForOf(evalInner, t.rhs.Context()); err != nil {
		return err
	}

	return nil
}

func (t *ForOf) Walk(fn WalkFunc) error {
  if err := t.ForInOf.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
