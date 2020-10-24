package js

import (
	"./values"

	"../context"
)

type ForIn struct {
	ForInOf
}

func NewForIn(varType VarType, lhs *VarExpression, rhs Expression, ctx context.Context) (*ForIn, error) {
	return &ForIn{newForInOf(varType, lhs, rhs, ctx)}, nil
}

func (t *ForIn) Dump(indent string) string {
	return t.ForInOf.dump(indent, "in")
}

func (t *ForIn) WriteStatement(indent string) string {
	return t.ForInOf.writeStatement(indent, "", "in")
}

func (t *ForIn) EvalStatement(stack values.Stack) error {
	rhsValue, err := t.rhs.EvalExpression(stack)
	if err != nil {
		return err
	}

	indexCtx := t.lhs.Context()

	// Objects always must be looped explicitely
	evalInner := func(v values.Value) error {
		subStack := NewBranchStack(stack)

		if t.varType == VAR {
			if err := stack.SetValue(t.lhs.GetVariable(), v, true,
				indexCtx); err != nil {
				return err
			}
		} else {
			if err := subStack.SetValue(t.lhs.GetVariable(), v, true,
				indexCtx); err != nil {
				return err
			}
		}

		return t.Block.EvalStatement(subStack)
	}

	if err := rhsValue.LoopForIn(evalInner, t.Context()); err != nil {
		return err
	}

	return nil
}

func (t *ForIn) Walk(fn WalkFunc) error { 
  if err := t.ForInOf.Walk(fn); err != nil {
    return err
  }

  return fn(t)
}
